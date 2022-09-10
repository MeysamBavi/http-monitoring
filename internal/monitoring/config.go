package monitoring

import "time"

type Config struct {
	RequestTimeout  time.Duration `config:"request_timeout"`
	NumberOfWorkers int           `config:"number_of_workers"`
}
