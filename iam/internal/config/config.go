package config

import (
	"github.com/alexander-kartavtsev/starship/iam/internal/config/env"
	"github.com/joho/godotenv"
	"os"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	Grpc     IamGrpcConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
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

	appConfig = &config{
		Logger:   loggerConfig,
		Grpc:     iamGrpcConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
		Session:  sessionConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
