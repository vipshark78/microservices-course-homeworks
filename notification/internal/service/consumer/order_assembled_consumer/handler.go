package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.DecodeOrderAssembled(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("OrderUUID", event.OrderUUID),
		zap.String("UserUUID", event.UserUUID),
		zap.String("EventUUID", event.EventUUID),
	)

	err = s.notificationService.SendOrderAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send order assembled notification", zap.Error(err))
		return err
	}

	return nil
}
