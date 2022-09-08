package monitor

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(logger *zap.Logger) {
	//Todo

	logger.Debug("starting the monitoring service")
}

func New(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use: "monitor",
		Short: "Starts the monitoring module for stored urls",
		Run: func(cmd *cobra.Command, args []string) {
			main(logger)
		},
	}
}