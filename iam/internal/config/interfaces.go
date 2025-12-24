package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OtelCollectorEndpoint() string
	ServiceName() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsDir() string
}
type IAMGRPCConfig interface {
	Address() string
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

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}
