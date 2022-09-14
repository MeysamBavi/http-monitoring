package request

import (
	"errors"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type URL struct {
	Url       string         `json:"url"`
	Threshold int            `json:"threshold"`
	Interval  model.Interval `json:"interval"`
}

func (url *URL) Validate() error {
	return validation.ValidateStruct(url,
		validation.Field(&url.Url, validation.Required, is.URL),
		validation.Field(&url.Threshold, validation.Required, validation.Min(5)),
		validation.Field(&url.Interval, validation.Required, validation.By(intervalMinRule)))
}

func intervalMinRule(value any) error {
	interval, ok := value.(model.Interval)
	if !ok {
		return errors.New("could not convert value to interval type")
	}

	if interval.Duration < time.Second*5 {
		return errors.New("interval must be at least 5s")
	}

	return nil
}
