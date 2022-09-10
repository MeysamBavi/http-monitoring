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
	// schedule writes on "in"
	// workers read from "in"
	// workers write on "out"
	// collect reads from "out"

	in := make(chan *Task, s.numOfWorkers)
	out := make(chan *Result, s.numOfWorkers)
	var wg sync.WaitGroup

	wg.Add(s.numOfWorkers)
	workers := make([]*Worker, s.numOfWorkers)

	for i := 0; i < s.numOfWorkers; i++ {
		workers[i] = NewWorker(s.requestTimeout, s.logger.Named(fmt.Sprintf("worker(%d)", i)))
		go workers[i].Work(&wg, in, out)
	}

	heap := s.initializeHeap()
	scheduleShutdown := make(chan int)
	go s.schedule(heap, in, scheduleShutdown)

	updateShutdown := make(chan int)
	updateDone := make(chan int)
	go s.update(heap, updateShutdown, updateDone)

	collectDone := make(chan int)
	go s.collect(out, collectDone)

	// wait for shutdown signal
	<-shutdown
	s.logger.Info("received shutdown signal")

	s.logger.Info("stopping update")
	updateShutdown <- 0
	<-updateDone

	s.logger.Info("stopping schedule")
	scheduleShutdown <- 0 // close "in"

	s.logger.Info("waiting for workers to finish")
	wg.Wait() // wait for workers to stop writing to "out"

	close(out) // close "out"
	s.logger.Info("waiting for collect to finish writing to db")
	<-collectDone // wait for collect to complete working
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

// reads from "out" and writes to database. sends signal on "done" when done
func (s *Scheduler) collect(out <-chan *Result, done chan<- int) {
	logger := s.logger.Named("collect")

	for r := range out {
		logger.Debug("saving this result to db", zap.Any("result", r))

		//Todo: save to database
	}

	done <- 0
}

// reads from db and updates heap
func (s *Scheduler) update(syncHeap *util.SyncHeap[*TimedURL], shutdown <-chan int, done chan<- int) {
	// Todo: listen for updates from database

	logger := s.logger.Named("update")
	i := uint32(0)

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
