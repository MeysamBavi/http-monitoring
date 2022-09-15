package migrate

import (
	"context"

	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/MeysamBavi/http-monitoring/internal/db"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	db, err := db.New(logger.Named("db"), cfg.Database)
	if err != nil {
		logger.Fatal("cannot create a db instance", zap.Error(err))
	}

	{
		idx, err := db.Collection(cfg.Database.UserCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		)

		if err != nil {
			logger.Fatal("cannot create username index", zap.Error(err))
		}

		logger.Info("database index created", zap.Any("index", idx))
	}

	{
		idx, err := db.Collection(cfg.Database.UrlCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.D{{Key: "user_id", Value: 1}},
			},
		)

		if err != nil {
			logger.Fatal("cannot create user id index", zap.Error(err))
		}

		logger.Info("database index created", zap.Any("index", idx))
	}

	{
		idx, err := db.Collection(cfg.Database.AlertCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.D{{Key: "url_id", Value: 1}},
			},
		)

		if err != nil {
			logger.Fatal("cannot create url id index", zap.Error(err))
		}

		logger.Info("database index created", zap.Any("index", idx))
	}
}

func New(cfg *config.Config, logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Creates indexes on mongodb collections",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, logger)
		},
	}
}
