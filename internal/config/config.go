package config

import (
	"runtime"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/db"
	"github.com/MeysamBavi/http-monitoring/internal/monitoring"
)

type Config struct {
	Debug      bool              `config:"debug"`
	Listen     string            `config:"listen"`
	Monitoring monitoring.Config `config:"monitoring"`
	Auth       auth.Config       `config:"auth"`
	Database   db.Config         `config:"database"`
}

func Default() Config {
	return Config{
		Debug:  true,
		Listen: ":1234",
		Monitoring: monitoring.Config{
			RequestTimeout:  10 * time.Second,
			NumberOfWorkers: runtime.NumCPU(),
		},
		Auth: auth.Config{
			SigningKey:  "veryBadSecret",
			ExpireAfter: 15 * time.Minute,
		},
		Database: db.Config{
			URI:                "mongodb://127.0.0.1:27017",
			DbName:             "httpm",
			ConnectionTimeout:  2 * time.Second,
			UserCollection:     "user",
			UrlCollection:      "url",
			AlertCollection:    "alert",
			UrlEventCollection: "url_event",
		},
	}
}
