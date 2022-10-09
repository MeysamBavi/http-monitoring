package request

import (
	"errors"
	"github.com/MeysamBavi/http-monitoring/internal/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Alert struct {
	UrlId string `param:"id" path:"id" description:"url id" required:"true"`
}

func (a *Alert) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.UrlId, validation.Required, validation.By(parsableId)),
	)
}

func parsableId(id any) error {
	idStr, ok := id.(string)
	if !ok {
		return errors.New("id is not a string")
	}

	if _, err := model.ParseId(idStr); err != nil {
		return errors.New("could not parse url id")
	}

	return nil
}

func (a *Alert) ParseUrlId() model.ID {
	id, err := model.ParseId(a.UrlId)
	if err != nil {
		panic(err)
	}
	return id
}
