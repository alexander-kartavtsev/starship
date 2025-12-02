package config

import (
	"github.com/joho/godotenv"

	"github.com/alexander-kartavtsev/starship/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger    LoggerConfig
	GRPC      GRPCConfig
	Mongo     MongoConfig
	IamClient GRPCConfig
	Tracing   TracingConfig
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

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	iamClientConfig, err := env.NewIamClientConfig()
	if err != nil {
		return err
	}

	tracingConfig, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:    loggerConfig,
		GRPC:      gRPCConfig,
		Mongo:     mongoConfig,
		IamClient: iamClientConfig,
		Tracing:   tracingConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
