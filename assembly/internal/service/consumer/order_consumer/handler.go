package order_consumer

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
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

	orderAssembledEvent := model.OrderAssembledEvent{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: 10,
	}

	err = s.orderAssembledProducer.ProduceOrderAssembled(ctx, orderAssembledEvent)
	if err != nil {
		logger.Error(ctx, "Failed to produce OrderAssembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "OrderAssembled event sent successfully",
		zap.String("OrderUUID", event.OrderUUID),
	)

	return nil
}
