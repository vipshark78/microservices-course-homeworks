package service

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUuid string, partUuids []string) (model.Order, error)
	OrderByUUID(ctx context.Context, orderUuid string) (model.Order, error)
	OrderCancel(ctx context.Context, orderUuid string) error
	OrderPay(ctx context.Context, orderUuid string, paymentMethod model.PaymentMethod) (string, error)
	OrderUpdateStatus(ctx context.Context, orderUuid string, orderStatus model.OrderStatus) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
