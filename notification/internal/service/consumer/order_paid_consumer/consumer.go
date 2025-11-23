package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/vipshark78/microservices-course-homeworks/notification/internal/converter/kafka"
	notificationService "github.com/vipshark78/microservices-course-homeworks/notification/internal/service"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

type service struct {
	orderPaidConsumer   kafka.Consumer
	orderPaidDecoder    kafkaConverter.OrderPaidDecoder
	notificationService notificationService.NotificationService
}

func NewService(consumer kafka.Consumer, decoder kafkaConverter.OrderPaidDecoder, notificationService notificationService.NotificationService) *service {
	return &service{
		orderPaidConsumer:   consumer,
		orderPaidDecoder:    decoder,
		notificationService: notificationService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
