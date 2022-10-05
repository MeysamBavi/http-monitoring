package api

import (
	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/MeysamBavi/http-monitoring/internal/db"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Setup(cfg *config.Config, logger *zap.Logger, app *echo.Echo) {
	s := getStore(cfg, logger)

	jh := getJwtHandler(cfg)

	registerAPIs(logger, s, jh, app)

	app.Debug = cfg.Debug
}

func registerAPIs(logger *zap.Logger, s store.Store, jh *auth.JwtHandler, app *echo.Echo) {
	logger = logger.Named("endpoint")

	uh := UserHandler{
		Logger:     logger.Named("user"),
		UserStore:  s.User(),
		JwtHandler: jh,
	}
	uh.Register(app.Group("/users"))

	urh := UrlHandler{
		Logger:     logger.Named("url"),
		UrlStore:   s.Url(),
		JwtHandler: jh,
	}
	urh.Register(app.Group("/urls"))

	ah := AlertHandler{
		Logger:     logger.Named("alert"),
		AlertStore: s.Alert(),
		JwtHandler: jh,
	}
	ah.Register(app.Group("/alerts"))
}

func getJwtHandler(cfg *config.Config) *auth.JwtHandler {
	return auth.NewJwtHandler(cfg.Auth)
}

func getStore(cfg *config.Config, logger *zap.Logger) store.Store {
	if cfg.InMemory {
		return store.NewInMemoryStore(logger.Named("in-memory"))
	}
	database, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("cannot create a database instance", zap.Error(err))
	}
	logger.Info("connected to mongo database", zap.String("name", database.Name()))

	return store.NewMongodbStore(
		database,
		cfg.Database,
		logger.Named("mongo"),
	)
}
