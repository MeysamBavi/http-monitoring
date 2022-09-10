package config

import (
	"runtime"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/monitoring"
)

type Config struct {
	Debug      bool              `config:"debug"`
	Listen     string            `config:"listen"`
	Monitoring monitoring.Config `config:"monitoring"`
}

func Default() Config {
	return Config{
		Debug:  true,
		Listen: ":1234",
		Monitoring: monitoring.Config{
			RequestTimeout:  10 * time.Second,
			NumberOfWorkers: runtime.NumCPU(),
		},
	}
}
