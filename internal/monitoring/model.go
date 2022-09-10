package monitoring

import "time"

type Task struct {
	URL    string
	UserId uint32
}

type Result struct {
	Task       *Task
	StatusCode int
	Body       string
}

type TimedURL struct {
	URL      string
	UserId   uint32
	Interval time.Duration
	callTime time.Time
	index    int
}

func NewTimedURL(URL string, UserId uint32, Interval time.Duration) *TimedURL {
	return &TimedURL{
		URL:      URL,
		UserId:   UserId,
		Interval: Interval,
	}
}
