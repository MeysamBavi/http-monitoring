package monitor

import (
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	//Todo

	logger.Debug("starting the monitoring service")
}

func New(cfg *config.Config, logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use: "monitor",
		Short: "Starts the monitoring module for stored urls",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, logger)
		},
	}
}