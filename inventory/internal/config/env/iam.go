package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamClientEnvConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type iamClientConfig struct {
	raw iamClientEnvConfig
}

func NewIamClientConfig() (*iamClientConfig, error) {
	var raw iamClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &iamClientConfig{raw: raw}, nil
}

func (cfg *iamClientConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *iamClientConfig) Port() string {
	return cfg.raw.Port
}
