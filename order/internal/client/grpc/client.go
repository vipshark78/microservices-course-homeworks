package grpc

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, payOrder model.PayOrder) (string, error)
}
