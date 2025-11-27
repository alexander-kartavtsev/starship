package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type IamGrpcConfig interface {
	Address() string
	Port() string
}

type PostgresConfig interface {
	Uri() string
	DatabaseName() string
	MigrationsDir() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	SessionTTL() time.Duration
}
