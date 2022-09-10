package monitoring

import (
	"fmt"
	"os"
	"sync"
	"time"

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

	heap := s.initialize()
	scheduleShutdown := make(chan int)
	go s.schedule(heap, in, scheduleShutdown)

	collectDone := make(chan int)
	go s.collect(out, collectDone)

	// wait for shutdown signal
	<-shutdown
	s.logger.Info("received shutdown signal")

	scheduleShutdown <- 0 // close "in"

	s.logger.Info("waiting for workers to finish their jobs")
	wg.Wait() // wait for workers to stop writing to "out"

	close(out) // close "out"
	s.logger.Info("waiting for collect to finish writing to db")
	<-collectDone // wait for collect to complete working
}

func (s *Scheduler) initialize() *Heap {
	//Todo: initialize heap from database

	//Todo: remove this
	l := []*TimedURL{
		NewTimeURL("https://httpbin.org/status/200", 1, 5*time.Second),
		NewTimeURL("https://httpbin.org/status/500", 1, 20*time.Second),
	}

	return NewHeap(l...)
}

// writes to "in" and closes it when shutdown signal is received
func (s *Scheduler) schedule(syncedHeap *Heap, in chan<- *Task, shutdown <-chan int) {
	logger := s.logger.Named("schedule")

	for {
		select {
		case <-shutdown:
			close(in)
			return

		default:
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
