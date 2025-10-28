package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	Uri() string
	DatabaseName() string
}
