package config

import (
	"time"

	"github.com/IBM/sarama"
)

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

type IamGRPCConfig interface {
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

type KafkaConfig interface {
	Brokers() []string
}

type OrderProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
