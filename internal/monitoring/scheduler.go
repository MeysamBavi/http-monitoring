package monitoring

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/MeysamBavi/http-monitoring/internal/util"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger         *zap.Logger
	numOfWorkers   int
	requestTimeout time.Duration
	dataStore      store.Store
}

func NewScheduler(logger *zap.Logger, numOfWorkers int, requestTimeout time.Duration, dataStore store.Store) *Scheduler {
	return &Scheduler{
		logger,
		numOfWorkers,
		requestTimeout,
		dataStore,
	}
}

func (s *Scheduler) Run(shutdown <-chan os.Signal) {

	scope := s.createScope(shutdown)
	s.startWorkers(scope)
	s.startModules(scope)
	s.waitForShutdown(scope)
}

// needed variables in run
type scope struct {
	// 'schedule' writes on "in"
	// workers read from "in"
	// workers write on "out"
	// 'collect' reads from "out"

	shutdown         <-chan os.Signal
	in               chan *Task
	out              chan *Result
	wg               sync.WaitGroup
	scheduleShutdown chan int
	updateShutdown   chan int
	updateDone       chan int
	collectDone      chan int
	syncHeap         *util.SyncHeap[*TimedURL]
}

func (s *Scheduler) createScope(shutdown <-chan os.Signal) *scope {
	return &scope{
		shutdown:         shutdown,
		in:               make(chan *Task, s.numOfWorkers),
		out:              make(chan *Result, s.numOfWorkers),
		scheduleShutdown: make(chan int),
		updateShutdown:   make(chan int),
		updateDone:       make(chan int),
		collectDone:      make(chan int),
		syncHeap:         nil,
	}
}

func (s *Scheduler) startWorkers(scope *scope) {
	scope.wg.Add(s.numOfWorkers)
	for i := 0; i < s.numOfWorkers; i++ {
		go NewWorker(s.requestTimeout, s.logger.Named(fmt.Sprintf("worker(%d)", i))).Work(&scope.wg, scope.in, scope.out)
	}
}

func (s *Scheduler) startModules(scope *scope) {
	scope.syncHeap = s.initializeHeap()
	s.logger.Info("starting modules")

	go s.schedule(scope.syncHeap, scope.in, scope.scheduleShutdown)
	go s.update(scope.syncHeap, scope.updateShutdown, scope.updateDone)
	go s.collect(scope.out, scope.collectDone)
}

func (s *Scheduler) waitForShutdown(scope *scope) {
	<-scope.shutdown
	s.logger.Info("received shutdown signal")

	s.logger.Info("stopping update")
	scope.updateShutdown <- 0
	<-scope.updateDone

	s.logger.Info("stopping schedule")
	scope.scheduleShutdown <- 0 // close "in"

	s.logger.Info("waiting for workers to finish")
	scope.wg.Wait() // wait for workers to stop writing to "out"

	close(scope.out) // close "out"
	s.logger.Info("waiting for collect to finish writing to db")
	<-scope.collectDone // wait for collect to complete working
}

func (s *Scheduler) initializeHeap() *util.SyncHeap[*TimedURL] {

	s.logger.Info("initializing urls heap")

	all := make([]*TimedURL, 0)
	err := s.dataStore.Url().ForAll(context.Background(), func(u model.URL) {
		t := NewTimedURL(u.Id, u.Url, u.UserId, u.Interval.Duration)
		all = append(all, t)
	})

	if err != nil {
		s.logger.Fatal("error reading all urls", zap.Error(err))
	}

	return util.NewSyncHeap[*TimedURL](NewHeap(all...))
}

// writes to "in" and closes it when shutdown signal is received
func (s *Scheduler) schedule(syncedHeap *util.SyncHeap[*TimedURL], in chan<- *Task, shutdown <-chan int) {
	logger := s.logger.Named("schedule")

	for {
		select {
		case <-shutdown:
			close(in)
			return

		default:
			if syncedHeap.Len() == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			earliestUrl := syncedHeap.Peek()
			if time.Now().Before(earliestUrl.callTime) {
				max := time.Millisecond * 100
				until := time.Until(earliestUrl.callTime)
				if until < max {
					time.Sleep(until)
				} else {
					time.Sleep(max)
				}
				continue
			}

			logger.Debug("sending this url to workers", zap.String("url", earliestUrl.URL))
			in <- &Task{
				UrlId:  earliestUrl.UrlId,
				URL:    earliestUrl.URL,
				UserId: earliestUrl.UserId,
			}

			earliestUrl.callTime = time.Now().Add(earliestUrl.Interval)
			syncedHeap.Fix(earliestUrl.index)
		}
	}
}

// reads from db and updates heap
func (s *Scheduler) update(syncHeap *util.SyncHeap[*TimedURL], shutdown <-chan int, done chan<- int) {
	logger := s.logger.Named("update")

	events, err := s.dataStore.Url().ListenForChanges(context.Background())
	if err != nil {
		logger.Fatal("error listening for changes", zap.Error(err))
	}

	for {
		select {
		case <-shutdown:
			done <- 0
			return
		case event, ok := <-events:
			if !ok {
				logger.Fatal("url events channel was closed unexpectedly")
			}
			logger.Debug("received event", zap.Any("event", event))
			if event.Operation == store.UrlChangeOperationInsert {
				logger.Debug("updating heap", zap.Any("event", event))
				syncHeap.Push(NewTimedURL(
					event.Url.Id,
					event.Url.Url,
					event.Url.UserId,
					event.Url.Interval.Duration,
				))
			}
		}
	}
}

// reads from "out" and writes to database. sends signal on "done" when done
func (s *Scheduler) collect(out <-chan *Result, done chan<- int) {
	logger := s.logger.Named("collect")

	for r := range out {
		logger.Debug("saving this result to db", zap.Any("result", r))

		var success, failure int
		if r.StatusCode >= 200 && r.StatusCode < 300 {
			success = 1
			failure = 0
		} else {
			success = 0
			failure = 1
		}

		url, stat, err := s.dataStore.Url().UpdateStat(context.Background(), r.Task.UserId, r.Task.UrlId, model.DayStat{
			Date:         model.Today(),
			SuccessCount: success,
			FailureCount: failure,
		})

		if err != nil {
			s.logger.Error("error updating stat", zap.Error(err), zap.Any("stat", stat))
			continue
		}

		// send alert if it has passed failure threshold
		if stat.FailureCount > 0 && stat.FailureCount%url.Threshold == 0 {
			err := s.dataStore.Alert().Add(context.Background(), &model.Alert{
				UserId:   url.UserId,
				UrlId:    url.Id,
				Url:      url.Url,
				IssuedAt: time.Now(),
			})

			if err != nil {
				s.logger.Error("error adding alert", zap.Error(err), zap.Any("url", url))
				continue
			}
		}
	}

	done <- 0
}
