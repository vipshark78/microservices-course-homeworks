package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	CollectorEndpoint string `env:"METRICS_OTEL_COLLECTOR_ENDPOINT,required"`
	CollectorInterval string `env:"METRICS_COLLECTOR_INTERVAL" envDefault:"10s"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &metricsConfig{raw: raw}, nil
}

func (cfg *metricsConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *metricsConfig) CollectorInterval() time.Duration {
	interval, err := time.ParseDuration(cfg.raw.CollectorInterval)
	if err != nil {
		return 10 * time.Second
	}
	return interval
}
