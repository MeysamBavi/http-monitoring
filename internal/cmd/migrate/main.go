package migrate

import (
	"context"
	"errors"

	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/MeysamBavi/http-monitoring/internal/db"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	database, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("cannot create a db instance", zap.Error(err))
	}
	Migrate(cfg, logger, database)
}

func Migrate(cfg *config.Config, logger *zap.Logger, db *mongo.Database) {
	{
		err := db.CreateCollection(
			context.Background(),
			cfg.Database.UrlEventCollection,
			options.CreateCollection().SetCapped(true).SetSizeInBytes(100_000*32),
		)

		if err != nil {
			var e mongo.CommandError
			if errors.As(err, &e) && e.Name == "NamespaceExists" {
				logger.Error("collection already exists", zap.Error(e))
			} else {
				logger.Fatal("cannot create url_event collection", zap.Error(err))
			}
		} else {
			logger.Info("url_event collection created")
		}
	}

	{
		idx, err := db.Collection(cfg.Database.UserCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.D{{Key: "username", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetCollation(&options.Collation{Locale: "en", Strength: 2}),
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
