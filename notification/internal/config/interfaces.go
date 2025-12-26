package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OtelCollectorEndpoint() string
	ServiceName() string
}

type KafkaConfig interface {
	Brokers() []string
}

type TelegramBotConfig interface {
	BotToken() string
	ChatID() int64
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
