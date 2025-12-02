package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/alexander-kartavtsev/starship/iam/internal/config/env"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	Grpc     IamGrpcConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
	Tracing  TracingConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	iamGrpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisConfig, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionConfig, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	tracingConfig, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:   loggerConfig,
		Grpc:     iamGrpcConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
		Session:  sessionConfig,
		Tracing:  tracingConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
