package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.DecodeOrderPaid(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("OrderUUID", event.OrderUUID),
		zap.String("UserUUID", event.UserUUID),
		zap.String("EventUUID", event.EventUUID),
		zap.String("TransactionUUID", event.TransactionUUID),
		zap.String("PaymentMethod", event.PaymentMethod),
	)

	err = s.notificationService.SendOrderPaidNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
		return err
	}

	return nil
}
