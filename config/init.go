package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type InitConfig struct {
	IocDryRun bool `env:"IOC_DRY_RUN" env-default:"false"`

	// Server configuration
	ServerPort int    `env:"SERVER_PORT" env-default:"8080"`
	ServerHost string `env:"SERVER_HOST" env-default:"localhost"`

	// Healthcheck configuration
	HealthcheckPort int `env:"HEALTHCHECK_PORT" env-default:"2802"`
}

func LoadInit() *InitConfig {
	var cfg InitConfig

	for _, file := range []string{".defaults.env", ".env"} {
		err := cleanenv.ReadConfig(file, &cfg)
		if err != nil {
			log.Panicf("failed to read config file %s: %v", file, err)
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Panicf("failed to parse config: %v", err)
	}

	return &cfg
}
