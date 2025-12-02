package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	"github.com/vipshark78/microservices-course-homeworks/notification/internal/client/http"
	telegramClient "github.com/vipshark78/microservices-course-homeworks/notification/internal/client/http/telegram"
	"github.com/vipshark78/microservices-course-homeworks/notification/internal/config"
	kafkaConverter "github.com/vipshark78/microservices-course-homeworks/notification/internal/converter/kafka"
	"github.com/vipshark78/microservices-course-homeworks/notification/internal/converter/kafka/decoder"
	"github.com/vipshark78/microservices-course-homeworks/notification/internal/service"
	orderAssembledConsumer "github.com/vipshark78/microservices-course-homeworks/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumer "github.com/vipshark78/microservices-course-homeworks/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/vipshark78/microservices-course-homeworks/notification/internal/service/telegram"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	wrappedKafka "github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka/consumer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
	kafkaMiddleware "github.com/vipshark78/microservices-course-homeworks/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderPaidConsumerService      service.OrderPaidConsumerService
	orderAssembledConsumerService service.OrderAssembledConsumerService
	orderPaidConsumerGroup        sarama.ConsumerGroup
	orderAssembledConsumerGroup   sarama.ConsumerGroup
	orderPaidConsumer             wrappedKafka.Consumer
	orderAssembledConsumer        wrappedKafka.Consumer
	orderPaidDecoder              kafkaConverter.OrderPaidDecoder
	orderAssembledDecoder         kafkaConverter.OrderAssembledDecoder
	notificationService           service.NotificationService
	telegramClient                http.TelegramClient
	telegramBot                   *bot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.NotificationService())
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerService() service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumer.NewService(d.OrderAssembledConsumer(), d.OrderAssembledDecoder(), d.NotificationService())
	}
	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create orderPaidConsumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka orderPaidConsumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{config.AppConfig().OrderAssembledConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create orderAssembledConsumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka orderAssembledConsumer group", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = consumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) NotificationService() service.NotificationService {
	if d.notificationService == nil {
		d.notificationService = telegramService.NewService(d.TelegramClient())
	}
	return d.notificationService
}

func (d *diContainer) TelegramClient() http.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot())
	}
	return d.telegramClient
}

func (d *diContainer) TelegramBot() *bot.Bot {
	if d.telegramBot == nil {
		bot, err := bot.New(config.AppConfig().TelegramBot.BotToken())
		if err != nil {
			panic(fmt.Errorf("failed to create telegram bot: %s", err.Error()))
		}
		d.telegramBot = bot
	}
	return d.telegramBot
}
