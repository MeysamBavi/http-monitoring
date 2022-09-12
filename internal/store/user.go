package store

import (
	"context"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type User interface {
	Get(context.Context, model.ID) (*model.User, error)
	Add(context.Context, *model.User) error
}
