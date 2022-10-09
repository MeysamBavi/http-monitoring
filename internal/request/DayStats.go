package request

import (
	"github.com/MeysamBavi/http-monitoring/internal/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DayStats struct {
	UrlId string `param:"id" path:"id" description:"url id" required:"true"`
	Day   *int   `query:"day" description:"day of the month (1-31)"`
	Month *int   `query:"month" description:"month number (1-12)"`
	Year  *int   `query:"year"`
}

func (d *DayStats) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.UrlId, validation.Required, validation.By(parsableId)),
		validation.Field(&d.Day, validation.Min(1), validation.Max(31)),
		validation.Field(&d.Month, validation.Min(1), validation.Max(12)),
	)
}

func (d *DayStats) DayFilter() func(date model.Date) bool {
	return func(date model.Date) bool {
		if d.Day != nil && *d.Day != date.Day {
			return false
		}

		if d.Month != nil && *d.Month != date.Month {
			return false
		}

		if d.Year != nil && *d.Year != date.Year {
			return false
		}

		return true
	}
}

func (d *DayStats) ParseUrlId() model.ID {
	id, err := model.ParseId(d.UrlId)
	if err != nil {
		panic(err)
	}
	return id
}
