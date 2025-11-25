package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type gRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type gRPCConfig struct {
	raw gRPCEnvConfig
}

func NewGRPCConfig() (*gRPCConfig, error) {
	var raw gRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &gRPCConfig{raw: raw}, nil
}

func (cfg *gRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *gRPCConfig) Port() string {
	return cfg.raw.Port
}
