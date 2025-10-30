package order

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository/converter"
)

func (r *repository) Read(ctx context.Context, uuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.orders[uuid]; !ok {
		return model.Order{}, model.ErrOrderNotFound
	}
	return converter.ModelToOrder(r.orders[uuid]), nil
}
