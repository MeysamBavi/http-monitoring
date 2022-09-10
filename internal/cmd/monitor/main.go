package monitor

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/MeysamBavi/http-monitoring/internal/monitoring"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	logger.Debug("starting the monitoring service")

	scheduler := monitoring.NewScheduler(
		logger.Named("scheduler"),
		cfg.Monitoring.NumberOfWorkers,
		cfg.Monitoring.RequestTimeout,
	)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	scheduler.Run(shutdown)
}

func New(cfg *config.Config, logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "monitor",
		Short: "Starts the monitoring module for stored urls",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, logger)
		},
	}
}
