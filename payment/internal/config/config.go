package config

import (
	"github.com/joho/godotenv"

	"github.com/alexander-kartavtsev/starship/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	GRPC   GRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	gRPCConfig, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerConfig,
		GRPC:   gRPCConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
