package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Minio MinioConfig
	Chat  ChatConfig
}

type AppConfig struct {
	Port         string   `env:"APP_PORT, default=7070"`
	Environment  string   `env:"APP_ENV, default=development"`
	MinioBuckets []string `env:"MINIO_BUCKETS, required"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST, required"`
	Port     int    `env:"DB_PORT, default=5432"`
	User     string `env:"DB_USER, required"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME, required"`
	SSLMode  string `env:"DB_SSLMODE, default=disable"`
}

type MinioConfig struct {
	Buckets   []string `env:"MINIO_BUCKETS, required"`
	Endpoint  string   `env:"MINIO_ENDPOINT, required"`
	AccessKey string   `env:"MINIO_ACCESS_KEY, required"`
	SecretKey string   `env:"MINIO_SECRET_KEY, required"`
	SSL       bool     `env:"MINIO_SSL, default=false"`
}

type ChatConfig struct {
	Provider      string `env:"CHAT_PROVIDER, default=github"`
	GitHubToken   string `env:"GITHUB_TOKEN"`
	OpenRouterKey string `env:"OPENROUTER_API_KEY"`
	Model         string `env:"CHAT_MODEL, default=openai/o4-mini"`
	SiteURL       string `env:"OPENROUTER_SITE_URL"`
	AppName       string `env:"OPENROUTER_APP_NAME"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}