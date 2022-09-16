package store

import (
	"context"
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

		return fmt.Errorf("document creation on user collection failed: %w", err)
	}

	doc.Id, _ = model.ParseId(r.InsertedID.(primitive.ObjectID).Hex()) // set the new id for caller

	return nil
}

func (m *MongodbUser) Get(ctx context.Context, id model.ID) (*model.User, error) {
	r := m.coll.FindOne(
		ctx,
		bson.M{"_id": id},
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
		return nil, NewNotFoundError("user", "username", username)
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

func (m *MongodbUrl) Add(context.Context, *model.URL) error {
	panic("unimplemented")
}

func (*MongodbUrl) ForAll(context.Context, func(model.URL)) error {
	panic("unimplemented")
}

func (*MongodbUrl) GetByUserId(context.Context, model.ID) ([]*model.URL, error) {
	panic("unimplemented")
}

func (*MongodbUrl) GetDayStat(ctx context.Context, userId model.ID, id model.ID, date model.Date) (model.DayStat, error) {
	panic("unimplemented")
}

func (*MongodbUrl) ListenForChanges(context.Context, chan<- UrlChangeEvent) error {
	panic("unimplemented")
}

func (*MongodbUrl) UpdateStat(ctx context.Context, userId model.ID, id model.ID, stat model.DayStat) (*model.URL, model.DayStat, error) {
	panic("unimplemented")
}

type MongodbAlert struct {
	coll *mongo.Collection
}

func (*MongodbAlert) Add(context.Context, *model.Alert) error {
	panic("unimplemented")
}

func (*MongodbAlert) GetByUrlId(context.Context, model.ID) ([]*model.Alert, error) {
	panic("unimplemented")
}
