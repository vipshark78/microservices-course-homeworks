package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/vipshark78/microservices-course-homeworks/notification/internal/converter/kafka"
	notificationService "github.com/vipshark78/microservices-course-homeworks/notification/internal/service"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	notificationService    notificationService.NotificationService
}

func NewService(consumer kafka.Consumer, decoder kafkaConverter.OrderAssembledDecoder, notificationService notificationService.NotificationService) *service {
	return &service{
		orderAssembledConsumer: consumer,
		orderAssembledDecoder:  decoder,
		notificationService:    notificationService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderAssembledConsumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
