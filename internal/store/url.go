package store

import (
	"context"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type Url interface {
	// blocks
	ListenForChanges(context.Context, chan<- UrlChangeEvent) error
	ForAll(context.Context, func(model.URL)) error
	GetByUserId(context.Context, model.ID) ([]*model.URL, error)
	GetDayStat(ctx context.Context, userId model.ID, id model.ID, date model.Date) (model.DayStat, error)
	Add(context.Context, *model.URL) error
	UpdateStat(ctx context.Context, userId model.ID, id model.ID, stat model.DayStat) (*model.URL, model.DayStat, error)
}

type UrlChangeEvent struct {
	Url       model.URL
	Operation int
}

const (
	UrlChangeOperationInsert = 1
	UrlChangeOperationUpdate = 2
	UrlChangeOperationDelete = 3
)
