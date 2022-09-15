package db

import "time"

type Config struct {
	URI               string        `config:"uri"`
	DbName            string        `config:"db_name"`
	ConnectionTimeout time.Duration `config:"connection_timeout"`
}
