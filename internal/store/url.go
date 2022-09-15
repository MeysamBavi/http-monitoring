package store

import (
	"context"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type Url interface {
	ForAll(context.Context, func(model.URL)) error
	GetByUserId(context.Context, model.ID) ([]*model.URL, error)
	GetDayStat(ctx context.Context, userId model.ID, id model.ID, date model.Date) (model.DayStat, error)
	Add(context.Context, *model.URL) error
	UpdateStat(ctx context.Context, userId model.ID, id model.ID, stat model.DayStat) (*model.URL, model.DayStat, error)
}
