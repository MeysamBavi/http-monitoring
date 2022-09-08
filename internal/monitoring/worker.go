package monitoring

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type Worker struct {
	ctx context.Context
	logger *zap.Logger
}

func NewWorker(ctx context.Context, logger *zap.Logger) *Worker {
	return &Worker{
		ctx, logger,
	}
}

func (w *Worker) Work(wg *sync.WaitGroup, in <-chan *Task, out chan<- *Result) {
	for t := range in {
		r, ok := w.Process(t)
		if ok {
			out <- r
		}
	}
	wg.Done()
}

func (w *Worker) Process(t *Task) (*Result, bool) {

	// creating http request
	req, errReq := http.NewRequestWithContext(w.ctx, http.MethodGet, t.URL, nil)

	if errReq != nil {
		w.logger.Error("error creating the request", zap.Error(errReq))
		return nil, false
	}

	// sending http request
	res, errRes := http.DefaultClient.Do(req)

	if errRes != nil {
		w.logger.Error("error sending the request", zap.Error(errRes))
		return nil, false
	}

	// reading response body
	var buf strings.Builder
	_, errResRead := io.Copy(&buf, res.Body)
	res.Body.Close()

	var responseBody string

	if errResRead != nil {
		w.logger.Error("error reading the response", zap.Error(errResRead))
		responseBody = ""
	} else {
		responseBody = buf.String()
	}

	return &Result{
		Task: t,
		StatusCode: res.StatusCode,
		Body: responseBody,
	}, true
}
