package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT,required"`
	PortExt  string `env:"EXTERNAL_MONGO_PORT,required"`
	Db       string `env:"MONGO_INITDB_DATABASE,required"`
	AuthDb   string `env:"MONGO_AUTH_DB,required"`
	User     string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &mongoConfig{raw: raw}, nil
}

func (cfg *mongoConfig) Uri() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.PortExt,
		cfg.raw.Db,
		cfg.raw.AuthDb,
	)
}

func (cfg *mongoConfig) DatabaseName() string {
	return cfg.raw.Db
}

func (cfg *mongoConfig) Port() string {
	return cfg.raw.PortExt
}
