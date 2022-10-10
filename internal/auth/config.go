package auth

import "time"

type Config struct {
	SigningKey  string        `config:"signing_key" json:"-"`
	ExpireAfter time.Duration `config:"expire_after"`
}
