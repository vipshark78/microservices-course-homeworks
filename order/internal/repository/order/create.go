package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository/converter"
	repomodel "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/model"
)

func (r *repository) Insert(ctx context.Context, userUuid string, partUuids []string, price float64) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	orderUUID := uuid.New()
	newOrder := repomodel.Order{
		UserUUID:        userUuid,
		OrderUUID:       orderUUID.String(),
		PartUuids:       partUuids,
		TotalPrice:      price,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          model.OrderStatusPENDINGPAYMENT,
	}

	r.orders[orderUUID.String()] = newOrder
	return converter.ModelToOrder(newOrder), nil
}
