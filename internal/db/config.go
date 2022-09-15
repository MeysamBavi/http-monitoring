package db

import "time"

type Config struct {
	URI               string        `config:"uri"`
	DbName            string        `config:"db_name"`
	UserCollection    string        `config:"user_collection"`
	UrlCollection     string        `config:"url_collection"`
	AlertCollection   string        `config:"alert_collection"`
	ConnectionTimeout time.Duration `config:"connection_timeout"`
}
