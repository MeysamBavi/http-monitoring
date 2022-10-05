package store

import (
	"context"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
)

type Url interface {
	ListenForChanges(context.Context) (<-chan UrlChangeEvent, error)
	ForAll(context.Context, func(model.URL)) error
	GetByUserId(context.Context, model.ID) ([]*model.URL, error)
	GetDayStats(ctx context.Context, userId model.ID, id model.ID, dateFilter func(model.Date) bool) ([]model.DayStat, error)
	Add(context.Context, *model.URL) error
	UpdateStat(ctx context.Context, userId model.ID, id model.ID, stat model.DayStat) (*model.URL, model.DayStat, error)
}

type UrlChangeEvent struct {
	Url       model.URL `bson:"fullDocument"`
	Operation string    `bson:"operationType"`
	Timestamp time.Time `bson:"timestamp"`
}

const (
	UrlChangeOperationInsert = "insert"
	UrlChangeOperationUpdate = "update"
	UrlChangeOperationDelete = "delete"
)
