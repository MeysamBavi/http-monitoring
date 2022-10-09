package serve

import (
	"github.com/MeysamBavi/http-monitoring/internal/api"
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	app := echo.New()

	api.Setup(cfg, logger, app)

	if err := app.Start(cfg.Listen); err != nil {
		logger.Fatal("cannot start the server", zap.Error(err))
	}
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
