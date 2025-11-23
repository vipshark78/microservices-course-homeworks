package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	InventoryGRPC          InventoryGRPCConfig
	PaymentGRPC            PaymentGRPCConfig
	Postgres               PostgresConfig
	OrderHTTP              OrderHTTPConfig
	Kafka                  KafkaConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	OrderPaidProducer      OrderPaidProducerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerConfig, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderPaidProducerConfig, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		InventoryGRPC:          inventoryGRPCCfg,
		PaymentGRPC:            paymentGRPCCfg,
		Postgres:               postgresCfg,
		OrderHTTP:              orderHTTPCfg,
		OrderAssembledConsumer: orderAssembledConsumerConfig,
		OrderPaidProducer:      orderPaidProducerConfig,
		Kafka:                  kafkaConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
