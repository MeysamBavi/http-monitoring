package monitoring

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/util"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger         *zap.Logger
	numOfWorkers   int
	requestTimeout time.Duration
}

func NewScheduler(logger *zap.Logger, numOfWorkers int, requestTimeout time.Duration) *Scheduler {
	return &Scheduler{
		logger,
		numOfWorkers,
		requestTimeout,
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
	//Todo: initialize heap from database

	return util.NewSyncHeap[*TimedURL](NewHeap())
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
				time.Sleep(time.Until(earliestUrl.callTime))
			}

			logger.Debug("sending this url to workers", zap.String("url", earliestUrl.URL))
			in <- &Task{
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
	// Todo: listen for updates from database

	logger := s.logger.Named("update")
	i := uint64(0)

	for {
		select {
		case <-shutdown:
			done <- 0
			return
		default:
			//Todo: remove this
			time.Sleep(10 * time.Second)
			logger.Debug("updating heap")
			syncHeap.Push(NewTimedURL("https://httpbin.org/status/206", i, 10*time.Second))
			i++
		}
	}
}

// reads from "out" and writes to database. sends signal on "done" when done
func (s *Scheduler) collect(out <-chan *Result, done chan<- int) {
	logger := s.logger.Named("collect")

	for r := range out {
		logger.Debug("saving this result to db", zap.Any("result", r))

		//Todo: save to database
	}

	done <- 0
}
