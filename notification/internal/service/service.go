package service

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/notification/internal/model"
)

type NotificationService interface {
	SendOrderPaidNotification(ctx context.Context, notification model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, notification model.OrderAssembledEvent) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}
