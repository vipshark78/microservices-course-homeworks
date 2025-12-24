package order_consumer

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/metrics"
	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	metrics.IncRequestCounter(ctx)
	startTime := time.Now()

	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("OrderUUID", event.OrderUUID),
		zap.String("TransactionUUID", event.TransactionUUID),
		zap.String("UserUUID", event.UserUUID),
		zap.String("EventUUID", event.EventUUID),
		zap.String("PaymentMethod", event.PaymentMethod),
	)

	logger.Info(ctx, "Waiting 10 seconds for assembly...",
		zap.String("OrderUUID", event.OrderUUID),
	)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
	}

	buildTimeSec, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		logger.Error(ctx, "Failed to generate random number", zap.Error(err))
		return err
	}

	orderAssembledEvent := model.OrderAssembledEvent{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildTimeSec.Int64(),
	}

	err = s.orderAssembledProducer.ProduceOrderAssembled(ctx, orderAssembledEvent)
	if err != nil {
		logger.Error(ctx, "Failed to produce OrderAssembled event", zap.Error(err))
		return err
	}

	duration := time.Since(startTime)
	durationSeconds := duration.Seconds()

	metrics.AssembledDuration.Record(ctx, durationSeconds,
		metric.WithAttributes(
			attribute.String("order_uuid", event.OrderUUID),
		),
	)
	metrics.AssembliesTotal.Add(ctx, 1)

	logger.Info(ctx, "OrderAssembled event sent successfully",
		zap.String("OrderUUID", event.OrderUUID),
	)

	return nil
}
