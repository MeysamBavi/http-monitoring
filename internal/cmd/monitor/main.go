package monitor

import (
	"github.com/MeysamBavi/http-monitoring/internal/cmd/migrate"
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

const (
	migrateFirstFlagName = "migrate-first"
)

func main(cfg *config.Config, logger *zap.Logger, migrateFirst bool) {

	logger.Debug("starting the monitoring service")

	var s store.Store
	if !cfg.InMemory {
		database, err := db.New(cfg.Database)
		if err != nil {
			logger.Fatal("cannot create a database instance", zap.Error(err))
		}
		logger.Info("connected to mongo database", zap.String("name", database.Name()))

		if migrateFirst {
			logger.Info("migrating database")
			migrate.Migrate(cfg, logger, database)
		}

		s = store.NewMongodbStore(
			database,
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

	logger.Info("running scheduler")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	scheduler.Run(shutdown)
}

func New(cfg *config.Config, logger *zap.Logger) *cobra.Command {
	migrateFirst := false
	command := &cobra.Command{
		Use:   "monitor",
		Short: "Starts the monitoring module for stored urls",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, logger, migrateFirst)
		},
	}
	command.Flags().BoolVarP(
		&migrateFirst,
		migrateFirstFlagName,
		"m",
		false,
		"Perform database migration before starting to monitor",
	)
	return command
}
