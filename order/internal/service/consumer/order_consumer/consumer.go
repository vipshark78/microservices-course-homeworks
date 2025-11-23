package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/vipshark78/microservices-course-homeworks/order/internal/converter/kafka"
	orderService "github.com/vipshark78/microservices-course-homeworks/order/internal/service"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderService           orderService.OrderService
}

func NewService(consumer kafka.Consumer, decoder kafkaConverter.OrderAssembledDecoder, orderService orderService.OrderService) *service {
	return &service{
		orderAssembledConsumer: consumer,
		orderAssembledDecoder:  decoder,
		orderService:           orderService,
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
