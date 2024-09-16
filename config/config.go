package config

type Config struct {
	IsLombok bool
}

func GetConfig() *Config {
	return &Config{
		IsLombok: true,
	}
}
