package env

import "github.com/caarlos0/env/v11"

type orderAssembledConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

func NewOrderAssembledConsumerConfig() (*ConsumerConfig, error) {
	var raw orderAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return NewConsumerConfig(raw.Topic, raw.GroupID), nil
}
