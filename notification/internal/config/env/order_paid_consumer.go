package env

import "github.com/caarlos0/env/v11"

type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

func NewOrderPaidConsumerConfig() (*ConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return NewConsumerConfig(raw.Topic, raw.GroupID), nil
}
