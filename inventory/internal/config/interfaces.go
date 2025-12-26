package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OtelCollectorEndpoint() string
	ServiceName() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type IAMConfig interface {
	Address() string
}
