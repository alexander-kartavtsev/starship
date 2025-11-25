package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
	Port() string
}

type MongoConfig interface {
	Uri() string
	DatabaseName() string
	Port() string
}
