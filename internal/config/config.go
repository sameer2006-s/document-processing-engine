package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	Environment string `env:"APP_ENV, default=development"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST, required"`
	Port     int    `env:"DB_PORT, default=5432"`
	User     string `env:"DB_USER, required"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME, required"`
	SSLMode  string `env:"DB_SSLMODE, default=disable"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}