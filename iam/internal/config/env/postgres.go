package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	PortExt       string `env:"EXTERNAL_POSTGRES_PORT,required"`
	Database      string `env:"POSTGRES_DB,required"`
	User          string `env:"POSTGRES_USER,required"`
	Password      string `env:"POSTGRES_PASSWORD,required"`
	MigrationsDir string `env:"MIGRATION_DIRECTORY,required"`
	SslMode       string `env:"POSTGRES_SSL_MODE,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) Uri() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.PortExt,
		cfg.raw.Database,
	)
}

func (cfg *postgresConfig) DatabaseName() string {
	return cfg.raw.Database
}

func (cfg *postgresConfig) MigrationsDir() string {
	return cfg.raw.MigrationsDir
}
