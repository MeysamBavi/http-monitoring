package serve

import (
	"github.com/MeysamBavi/http-monitoring/internal/api"
	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/MeysamBavi/http-monitoring/internal/db"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	app := echo.New()

	var s store.Store
	{
		logger := logger.Named("store")
		db, err := db.New(cfg.Database)
		if err != nil {
			logger.Fatal("cannot create a db instance", zap.Error(err))
		}
		logger.Info("connected to mongo db", zap.String("name", db.Name()))

		s = store.NewMongodbStore(
			logger.Named("mongo"),
			db,
			db.Collection(cfg.Database.UserCollection),
			db.Collection(cfg.Database.UrlCollection),
			db.Collection(cfg.Database.AlertCollection),
			db.Collection(cfg.Database.UrlEventCollection),
		)
	}

	var jh *auth.JwtHandler
	{
		jh = auth.NewJwtHandler(cfg.Auth)
	}

	{
		logger := logger.Named("endpoint")

		uh := api.UserHandler{
			Logger:     logger.Named("user"),
			UserStore:  s.User(),
			JwtHandler: jh,
		}
		uh.Register(app.Group("/user"))

		urh := api.UrlHandler{
			Logger:     logger.Named("url"),
			UrlStore:   s.Url(),
			JwtHandler: jh,
		}
		urh.Register(app.Group("/url"))

		ah := api.AlertHandler{
			Logger:     logger.Named("alert"),
			AlertStore: s.Alert(),
			JwtHandler: jh,
		}
		ah.Register(app.Group("/alert"))
	}

	app.Debug = cfg.Debug

	app.Start(cfg.Listen)
}

func New(cfg *config.Config, logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Runs the http server",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, logger)
		},
	}
}
