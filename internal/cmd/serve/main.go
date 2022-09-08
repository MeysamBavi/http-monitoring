package serve

import (
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	//Todo

	app := echo.New()
	app.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	app.Debug = cfg.Debug

	app.Start(cfg.Listen)
}


func New(cfg *config.Config, logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use: "serve",
		Short: "Runs the http server",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, logger);
		},
	}
}