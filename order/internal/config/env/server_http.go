package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type serverHttpEnvConfig struct {
	Host            string `env:"HTTP_HOST,required"`
	Port            string `env:"HTTP_PORT,required"`
	TimeoutRead     int    `env:"HTTP_READ_TIMEOUT"`
	TimeoutShutdown int    `env:"HTTP_SHUTDOWN_TIMEOUT"`
}

type serverHttpConfig struct {
	raw serverHttpEnvConfig
}

func NewServerHttpConfig() (*serverHttpConfig, error) {
	var raw serverHttpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &serverHttpConfig{raw: raw}, nil
}

func (cfg serverHttpConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
	// return fmt.Sprintf("%s:%s", cfg.raw.Host, cfg.raw.Port)
}

func (cfg serverHttpConfig) Port() string {
	return cfg.raw.Port
}

func (cfg serverHttpConfig) TimeoutRead() time.Duration {
	timeout := cfg.raw.TimeoutRead
	if int(timeout) == 0 {
		timeout = 5
	}
	return time.Duration(timeout) * time.Second
}

func (cfg serverHttpConfig) TimeoutShutdown() time.Duration {
	timeout := cfg.raw.TimeoutShutdown
	if int(timeout) == 0 {
		timeout = 10
	}
	return time.Duration(timeout) * time.Second
}
