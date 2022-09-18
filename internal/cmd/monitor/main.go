package monitor

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/MeysamBavi/http-monitoring/internal/db"
	"github.com/MeysamBavi/http-monitoring/internal/monitoring"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(cfg *config.Config, logger *zap.Logger) {
	logger.Debug("starting the monitoring service")

	var s store.Store
	if !cfg.InMemory {
		db, err := db.New(cfg.Database)
		if err != nil {
			logger.Fatal("cannot create a db instance", zap.Error(err))
		}
		logger.Info("connected to mongo db", zap.String("name", db.Name()))

		s = store.NewMongodbStore(
			db,
			cfg.Database,
			logger.Named("mongo"),
		)
	} else {
		s = store.NewInMemoryStore(logger.Named("in-memory"))
	}

	scheduler := monitoring.NewScheduler(
		logger.Named("scheduler"),
		cfg.Monitoring.NumberOfWorkers,
		cfg.Monitoring.RequestTimeout,
		s,
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
