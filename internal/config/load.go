package config

import (
	"github.com/knadh/koanf/providers/env"
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

const (
	tag       = "config"
	delimiter = "."
	prefix    = "HTTPM_"
	separator = "__"
)

func Load() *Config {
	k := koanf.New(delimiter)

	{
		err := k.Load(structs.Provider(Default(), tag), nil)
		if err != nil {
			log.Fatalf("could not load default config: %s", err)
		}
	}

	{
		err := k.Load(file.Provider("config.json"), json.Parser())
		if err != nil {
			log.Printf("could not load json config: %s\n", err)
		}
	}

	{
		err := k.Load(env.Provider(prefix, delimiter, envCallBack), nil)
		if err != nil {
			log.Printf("could not load env variables for config: %s\n", err)
		}
	}

	var instance Config
	err := k.UnmarshalWithConf("", &instance, koanf.UnmarshalConf{
		Tag: tag,
	})

	if err != nil {
		log.Fatalf("could not unmarshal config: %s\n", err)
	}

	return &instance
}

func envCallBack(s string) string {
	base := strings.ToLower(strings.TrimPrefix(s, prefix))

	return strings.ReplaceAll(base, separator, delimiter)
}
