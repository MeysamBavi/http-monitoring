package api

import (
	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/MeysamBavi/http-monitoring/internal/db"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func Setup(cfg *config.Config, logger *zap.Logger, app *echo.Echo) {
	s := getStore(cfg, logger)

	jh := getJwtHandler(cfg)

	app.Use(newLoggerMiddleware(logger))
	app.Use(middleware.RequestID())
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

func newLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			logger.Info("response sent",
				zap.String("id", c.Response().Header().Get(echo.HeaderXRequestID)),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", c.Request().Host),
				zap.String("method", c.Request().Method),
				zap.String("URI", c.Request().RequestURI),
				zap.String("user_agent", c.Request().UserAgent()),
				zap.Int("status", c.Response().Status),
				zap.Error(err),
				zap.Int64("bytes_in", c.Request().ContentLength),
				zap.Int64("bytes_out", c.Response().Size),
			)
			return err
		}
	}
}
