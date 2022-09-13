package monitoring

import (
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type Task struct {
	URL    string
	UserId model.ID
}

type Result struct {
	Task       *Task
	StatusCode int
	Body       string
}

type TimedURL struct {
	URL      string
	UserId   model.ID
	Interval time.Duration
	callTime time.Time
	index    int
}

func NewTimedURL(URL string, UserId uint64, Interval time.Duration) *TimedURL {
	return &TimedURL{
		URL:      URL,
		UserId:   model.ID(UserId),
		Interval: Interval,
	}
}
