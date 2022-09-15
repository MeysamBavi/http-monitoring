package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func New(logger *zap.Logger, cfg Config) (*mongo.Database, error) {
	options := options.Client()
	options.ApplyURI(cfg.URI)
	options.SetConnectTimeout(cfg.ConnectionTimeout)

	client, err := mongo.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	// connect to db
	if err := client.Connect(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// ping db
	{
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
		defer cancel()

		if err := client.Ping(ctx, nil); err != nil {
			return nil, fmt.Errorf("failed to ping mongodb: %w", err)
		}
	}

	return client.Database(cfg.DbName), nil
}
