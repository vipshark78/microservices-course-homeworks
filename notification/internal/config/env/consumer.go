package env

import "github.com/IBM/sarama"

type ConsumerConfig struct {
	topic   string
	groupID string
}

func NewConsumerConfig(topic, groupID string) *ConsumerConfig {
	return &ConsumerConfig{
		topic:   topic,
		groupID: groupID,
	}
}

func (cfg *ConsumerConfig) Topic() string {
	return cfg.topic
}

func (cfg *ConsumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *ConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
