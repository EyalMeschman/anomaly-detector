package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type InitConfig struct {
	IocDryRun bool `env:"IOC_DRY_RUN" env-default:"false"`

	// Server configuration
	ServerPort string `env:"SERVER_PORT" env-default:"8080"`
	ServerHost string `env:"SERVER_HOST" env-default:"localhost"`
}

func LoadInit() *InitConfig {
	var cfg InitConfig

	for _, file := range []string{".env.defaults", ".env"} {
		_ = cleanenv.ReadConfig(file, &cfg)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Panicf("failed to parse config: %v", err)
	}

	return &cfg
}
