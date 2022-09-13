package model

import "time"

type URL struct {
	Id        ID
	UserId    ID
	Url       string
	Threshold int
	Interval  time.Duration
	DayStats  []*DayStat
}

type DayStat struct {
	UrlId        ID
	Date         time.Time
	SuccessCount int
	FailureCount int
}
