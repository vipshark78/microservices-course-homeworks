package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/vipshark78/microservices-course-homeworks/assembly/internal/converter/kafka"
	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/metrics"
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
		metrics.IncResponseCounter(ctx, "error", "order_paid_consumer.RunConsumer")
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}
	metrics.IncResponseCounter(ctx, "success", "order_paid_consumer.RunConsumer")
	return nil
}
