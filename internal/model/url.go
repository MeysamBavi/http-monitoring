package model

import (
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type URL struct {
	Id        ID         `json:"id" bson:"_id"`
	UserId    ID         `json:"-" bson:"user_id"`
	Url       string     `json:"url" bson:"url"`
	Threshold int        `json:"threshold" bson:"threshold"`
	Interval  Interval   `json:"interval" bson:"interval"`
	DayStats  []*DayStat `json:"-" bson:"day_stats"`
}

func (u *URL) NoId() bson.M {
	return bson.M{
		"user_id":   u.UserId,
		"url":       u.Url,
		"threshold": u.Threshold,
		"interval":  u.Interval,
		"day_stats": u.DayStats,
	}
}

type DayStat struct {
	Date         Date `json:"date" bson:"date"`
	SuccessCount int  `json:"success_count" bson:"success_count"`
	FailureCount int  `json:"failure_count" bson:"failure_count"`
}

type Date struct {
	Day   int `json:"day" bson:"day"`
	Month int `json:"month" bson:"month"`
	Year  int `json:"year" bson:"year"`
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
