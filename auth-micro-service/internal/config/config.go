package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"APP_NAME"`
	Debug       bool   `envconfig:"APP_DEBUG"`
	GRPCPort    string `envconfig:"APP_GRPC_ADDRESS"`
	Secret      string `envconfig:"APP_SECRET"`
	Postgres    PostgresConfig
}

type PostgresConfig struct {
	Host          string `envconfig:"POSTGRES_HOST"`
	Port          string `envconfig:"POSTGRES_PORT"`
	User          string `envconfig:"POSTGRES_USER"`
	Password      string `envconfig:"POSTGRES_PASSWORD"`
	Database      string `envconfig:"POSTGRES_DB"`
	SSLMode       string `envconfig:"POSTGRES_SSLMODE"`
	MaxConnection int    `envconfig:"POSTGRES_MAX_CONNECTION" default:"10"`
	MinConnection int    `envconfig:"POSTGRES_MIN_CONNECTION" default:"0"`
}

func FromEnv() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("error while parse env config | %w", err)
	}

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database,
		c.Postgres.SSLMode,
	)
}
