package serve

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(logger *zap.Logger) {
	//Todo

	app := echo.New()
	app.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	app.Start(":1234")
}


// Todo Config: add config param to main and new
func New(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use: "serve",
		Short: "Runs the http server",
		Run: func(cmd *cobra.Command, args []string) {
			main(logger);
		},
	}
}