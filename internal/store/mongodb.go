package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MongodbStore struct {
	db     *mongo.Database
	logger *zap.Logger
	user   *MongodbUser
	url    *MongodbUrl
	alert  *MongodbAlert
}

func NewMongodbStore(logger *zap.Logger, db *mongo.Database, user, url, alert *mongo.Collection) Store {
	return &MongodbStore{
		db:     db,
		logger: logger,
		user:   &MongodbUser{user},
		url:    &MongodbUrl{url},
		alert:  &MongodbAlert{alert},
	}
}

func (s *MongodbStore) User() User {
	return s.user
}

func (s *MongodbStore) Url() Url {
	return s.url
}

func (s *MongodbStore) Alert() Alert {
	return s.alert
}

type MongodbUser struct {
	coll *mongo.Collection
}

func (m *MongodbUser) Add(ctx context.Context, doc *model.User) error {
	r, err := m.coll.InsertOne(
		ctx,
		doc.NoId(), // pass the document without _id field to generate new id
	)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return NewDuplicateError("user", "username", doc.Username)
		}

		return fmt.Errorf("error creating document: %w", err)
	}

	doc.Id = model.ParseIdFromObjectId(r.InsertedID.(primitive.ObjectID)) // set the new id for caller

	return nil
}

func (m *MongodbUser) Get(ctx context.Context, id model.ID) (*model.User, error) {
	r := m.coll.FindOne(
		ctx,
		bson.M{"_id": id.ObjectId()},
	)

	if r.Err() != nil {
		return nil, NewNotFoundError("user", "id", id)
	}

	var user model.User
	if err := r.Decode(&user); err != nil {
		return nil, fmt.Errorf("could not decode result into user: %w", err)
	}

	return &user, nil
}

func (m *MongodbUser) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	r := m.coll.FindOne(
		ctx,
		bson.M{"username": username},
	)

	if r.Err() != nil {
		if r.Err() == mongo.ErrNoDocuments {
			return nil, NewNotFoundError("user", "username", username)
		}

		return nil, fmt.Errorf("error getting url: %w", r.Err())
	}

	var user model.User
	if err := r.Decode(&user); err != nil {
		return nil, fmt.Errorf("could not decode result into user: %w", err)
	}

	return &user, nil
}

type MongodbUrl struct {
	coll *mongo.Collection
}

func (m *MongodbUrl) Add(ctx context.Context, doc *model.URL) error {
	doc.DayStats = make([]*model.DayStat, 0)

	r, err := m.coll.InsertOne(
		ctx,
		doc.NoId(), // pass the document without _id field to generate new id
	)

	if err != nil {
		return fmt.Errorf("error creating document: %w", err)
	}

	doc.Id = model.ParseIdFromObjectId(r.InsertedID.(primitive.ObjectID)) // set the new id for caller

	return nil
}

func (m *MongodbUrl) ForAll(ctx context.Context, action func(model.URL)) error {
	cursor, err := m.coll.Find(ctx, bson.D{})

	if err != nil {
		return fmt.Errorf("error reading all documents: %w", err)
	}

	for cursor.Next(ctx) {
		if cursor.Err() != nil {
			return fmt.Errorf("error reading from cursor: %w", err)
		}

		var url model.URL
		if err := cursor.Decode(&url); err != nil {
			return fmt.Errorf("error decoding current cursor value to url: %w", err)
		}

		action(url)
	}

	return nil
}

func (m *MongodbUrl) GetByUserId(ctx context.Context, id model.ID) ([]*model.URL, error) {
	cursor, err := m.coll.Find(
		ctx,
		bson.M{
			"user_id": id,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error reading from url collection: %w", err)
	}

	all := make([]*model.URL, 0)
	if err := cursor.All(ctx, &all); err != nil {
		return nil, fmt.Errorf("error decoding all results to url: %w", err)
	}

	return all, nil
}

func (m *MongodbUrl) GetDayStat(ctx context.Context, userId model.ID, id model.ID, date model.Date) (model.DayStat, error) {
	r := m.coll.FindOne(
		ctx,
		bson.M{
			"_id":     id.ObjectId(),
			"user_id": userId,
		},
	)

	if r.Err() != nil {
		if r.Err() == mongo.ErrNoDocuments {
			return model.DayStat{}, NotFoundError("found no url matching the parameters")
		}

		return model.DayStat{}, fmt.Errorf("error getting url: %w", r.Err())
	}

	var url model.URL
	if err := r.Decode(&url); err != nil {
		return model.DayStat{}, fmt.Errorf("could not decode result into url: %w", err)
	}

	for _, stat := range url.DayStats {
		if stat.Date == date {
			return *stat, nil
		}
	}

	return model.DayStat{}, NewNotFoundError("dayStat", "date", date)
}

func (m *MongodbUrl) ListenForChanges(ctx context.Context, events chan<- UrlChangeEvent) error {
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"operationType", bson.D{{
					"$in", bson.A{UrlChangeOperationInsert, UrlChangeOperationUpdate, UrlChangeOperationDelete}}},
				}},
			},
			{"$project", bson.D{
				{"fullDocument", 1},
				{"operationType", 1},
			}},
		},
	}

	cs, err := m.coll.Watch(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("error listening to changes: %w", err)
	}

	for cs.Next(ctx) {
		var event UrlChangeEvent
		if err := cs.Decode(&event); err != nil {
			return fmt.Errorf("error decoding change event: %w", err)
		}

		events <- event
	}

	return nil
}

func (m *MongodbUrl) UpdateStat(ctx context.Context, userId model.ID, id model.ID, stat model.DayStat) (*model.URL, model.DayStat, error) {
	r := m.coll.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id":     id.ObjectId(),
			"user_id": userId,
			"day_stats": bson.M{
				"$elemMatch": bson.M{
					"date": stat.Date,
				},
			},
		},
		bson.M{
			"$inc": bson.M{
				"day_stats.$.success_count": stat.SuccessCount,
				"day_stats.$.failure_count": stat.FailureCount,
			},
		},
	)

	if r.Err() != nil {
		if r.Err() != mongo.ErrNoDocuments {
			return nil, model.DayStat{}, fmt.Errorf("error updating url stat: %w", r.Err())
		}

		// no stat found, create a new one
		r = m.coll.FindOneAndUpdate(
			ctx,
			bson.M{
				"_id":     id.ObjectId(),
				"user_id": userId,
			},
			bson.M{
				"$push": bson.M{
					"day_stats": stat,
				},
			},
		)

		if r.Err() != nil {
			if r.Err() == mongo.ErrNoDocuments {
				return nil, model.DayStat{}, NotFoundError("found no url matching the parameters")
			}
			return nil, model.DayStat{}, fmt.Errorf("error updating url stat: %w", r.Err())
		}

		var url model.URL
		if err := r.Decode(&url); err != nil {
			return nil, model.DayStat{}, fmt.Errorf("could not decode result into url: %w", err)
		}

		return &url, stat, nil
	}

	var url model.URL
	if err := r.Decode(&url); err != nil {
		return nil, model.DayStat{}, fmt.Errorf("could not decode result into url: %w", err)
	}

	for _, s := range url.DayStats {
		if s.Date == stat.Date {
			return &url, *s, nil
		}
	}

	return nil, model.DayStat{}, errors.New("could not find updated stat")
}

type MongodbAlert struct {
	coll *mongo.Collection
}

func (m *MongodbAlert) Add(ctx context.Context, alert *model.Alert) error {
	r, err := m.coll.InsertOne(ctx, alert.NoId())
	if err != nil {
		return fmt.Errorf("error inserting alert: %w", err)
	}

	alert.Id = model.ParseIdFromObjectId(r.InsertedID.(primitive.ObjectID))

	return nil
}

func (m *MongodbAlert) GetByUrlId(ctx context.Context, id model.ID) ([]*model.Alert, error) {
	cursor, err := m.coll.Find(
		ctx,
		bson.M{
			"url_id": id,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error reading from alert collection: %w", err)
	}

	all := make([]*model.Alert, 0)
	if err := cursor.All(ctx, &all); err != nil {
		return nil, fmt.Errorf("error decoding all results to alert: %w", err)
	}

	return all, nil
}
