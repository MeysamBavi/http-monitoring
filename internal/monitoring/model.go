package monitoring

import (
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type Task struct {
	UrlId  model.ID
	URL    string
	UserId model.ID
}

type Result struct {
	Task       *Task
	StatusCode int
	Body       string
}

type TimedURL struct {
	UrlId    model.ID
	URL      string
	UserId   model.ID
	Interval time.Duration
	callTime time.Time
	index    int
}

func NewTimedURL(UrlId model.ID, URL string, UserId model.ID, Interval time.Duration) *TimedURL {
	return &TimedURL{
		UrlId:    UrlId,
		URL:      URL,
		UserId:   UserId,
		Interval: Interval,
		callTime: time.Now().Add(Interval), // for preventing starting the first call immediately
	}
}
