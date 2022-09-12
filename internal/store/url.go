package store

import (
	"context"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type Url interface {
	GetByUserId(context.Context, model.ID) ([]*model.URL, error)
	Add(context.Context, *model.URL) error
	AddStat(ctx context.Context, userId model.ID, id model.ID, stat *model.DayStat) error
}
