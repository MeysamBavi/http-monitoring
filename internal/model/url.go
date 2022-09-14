package model

import (
	"encoding/json"
	"errors"
	"time"
)

type URL struct {
	Id        ID         `json:"id"`
	UserId    ID         `json:"-"`
	Url       string     `json:"url"`
	Threshold int        `json:"threshold"`
	Interval  Interval   `json:"interval"`
	DayStats  []*DayStat `json:"-"`
}

type DayStat struct {
	UrlId        ID   `json:"-"`
	Date         Date `json:"date"`
	SuccessCount int  `json:"success_count"`
	FailureCount int  `json:"failure_count"`
}

type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

func ParseDate(s string) (Date, error) {
	date, err := time.Parse("2006/01/02", s)
	return Date{
		Day:   date.Day(),
		Month: int(date.Month()),
		Year:  date.Year(),
	}, err
}

func Today() Date {
	date := time.Now()
	return Date{
		Day:   date.Day(),
		Month: int(date.Month()),
		Year:  date.Year(),
	}
}

type Interval struct {
	time.Duration
}

func (i *Interval) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *Interval) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		i.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		i.Duration, err = time.ParseDuration(value)
		return err
	default:
		return errors.New("invalid interval")
	}
}
