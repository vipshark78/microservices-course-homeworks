package service

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUuid string, partUuids []string) (model.Order, error)
	OrderByUUID(ctx context.Context, orderUuid string) (model.Order, error)
	OrderCancel(ctx context.Context, orderUuid string) error
	OrderPay(ctx context.Context, orderUuid, paymentMethod string) (string, error)
}
