package cmd

import (
	"log"
	"os"

	"github.com/MeysamBavi/http-monitoring/internal/cmd/monitor"
	"github.com/MeysamBavi/http-monitoring/internal/cmd/serve"
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Execute() {
	//Todo config: select logger type based
	config := config.Load()
	log.Println(*config)

	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatal(err)
	}

	root := cobra.Command{
		Use: "httpm",
		Short: "http monitoring service - summer 2022",
	}

	root.AddCommand(serve.New(logger))
	root.AddCommand(monitor.New(logger))

	if err := root.Execute(); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
}
