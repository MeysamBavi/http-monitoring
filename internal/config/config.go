package config

type Config struct {
	Debug  bool   `config:"debug"`
	Listen string `config:"listen"`
}

func Default() Config {
	return Config{
		Debug:  true,
		Listen: ":1234",
	}
}
