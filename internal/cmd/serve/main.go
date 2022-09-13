package serve

import (
	"github.com/MeysamBavi/http-monitoring/internal/api"
	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/config"
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
		s = store.NewInMemoryStore(logger.Named("memory"))
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
