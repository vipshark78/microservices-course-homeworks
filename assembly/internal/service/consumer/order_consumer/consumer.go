package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/vipshark78/microservices-course-homeworks/assembly/internal/converter/kafka"
	assemblyService "github.com/vipshark78/microservices-course-homeworks/assembly/internal/service"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

type service struct {
	orderPaidConsumer      kafka.Consumer
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder
	orderAssembledProducer assemblyService.ProducerService
}

func NewService(consumer kafka.Consumer, decoder kafkaConverter.OrderPaidDecoder, producer assemblyService.ProducerService) *service {
	return &service{
		orderPaidConsumer:      consumer,
		orderPaidDecoder:       decoder,
		orderAssembledProducer: producer,
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
