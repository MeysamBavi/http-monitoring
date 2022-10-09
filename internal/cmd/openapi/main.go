package openapi

import (
	"github.com/MeysamBavi/http-monitoring/internal/api/apidoc"
	"github.com/MeysamBavi/http-monitoring/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

const (
	outputFlagName = "output"
)

func main(_ *config.Config, logger *zap.Logger, output string) {
	logger.Debug("generating OpenAPI 3.0 specification")

	file, err := os.Create(output)
	if err != nil {
		logger.Fatal("cannot create output file", zap.Error(err))
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.Fatal("cannot close output file", zap.Error(err))
		}
	}()

	spec, err := apidoc.NewDocGenerator(logger).OpenAPISpecAsYaml()
	if err != nil {
		logger.Fatal("cannot generate OpenAPI 3.0 specification", zap.Error(err))
	}

	if _, err := file.Write(spec); err != nil {
		logger.Fatal("cannot write specification to file", zap.Error(err))
	}

	logger.Info("OpenAPI 3.0 generated", zap.String("file", file.Name()))
}

func New(cfg *config.Config, logger *zap.Logger) *cobra.Command {
	output := ""

	command := &cobra.Command{
		Use:   "openapi",
		Short: "Generates OpenAPI 3.0 specification",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg, logger, output)
		},
	}

	command.Flags().StringVarP(
		&output,
		outputFlagName,
		"o",
		"openapi/httpm.yaml",
		"output file name",
	)

	return command
}
