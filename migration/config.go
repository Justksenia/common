package migration

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Postgres struct {
	DSN string `envconfig:"POSTGRES_DSN" required:"true"`
}

type Config struct {
	Postgres Postgres
}

func NewConfigFromEnv() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("process config from env: %w", err)
	}

	return &cfg, nil
}
