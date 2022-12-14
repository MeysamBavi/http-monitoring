package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MeysamBavi/http-monitoring/internal/cmd/migrate"
	"github.com/MeysamBavi/http-monitoring/internal/cmd/monitor"
	"github.com/MeysamBavi/http-monitoring/internal/cmd/openapi"
	"github.com/MeysamBavi/http-monitoring/internal/cmd/serve"
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() {
	cfg := config.Load()

	{
		cfgJson, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("config:", string(cfgJson))
	}

	var (
		logger *zap.Logger
		err    error
	)

	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatal(err)
	}

	root := cobra.Command{
		Use:   "httpm",
		Short: "http monitoring service - summer 2022",
	}

	root.AddCommand(serve.New(cfg, logger))
	root.AddCommand(monitor.New(cfg, logger))
	root.AddCommand(migrate.New(cfg, logger))
	root.AddCommand(openapi.New(cfg, logger))

	if err := root.Execute(); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
}
