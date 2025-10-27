package repository

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

type OrderRepository interface {
	Insert(ctx context.Context, userUuid string, partUuids []string, price float64) (model.Order, error)
	Read(ctx context.Context, uuid string) (model.Order, error)
	Update(order model.Order) error
}
