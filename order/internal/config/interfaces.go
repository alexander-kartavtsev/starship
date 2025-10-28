package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	Uri() string
	DatabaseName() string
	MigrationsDir() string
}

type ServerHTTPConfig interface {
	Address() string
	Port() string
	TimeoutRead() time.Duration
	TimeoutShutdown() time.Duration
}
