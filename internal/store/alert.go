package store

import (
	"context"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type Alert interface {
	GetByUrlId(context.Context, model.ID) ([]*model.Alert, error)
	Add(context.Context, *model.Alert) error
}
