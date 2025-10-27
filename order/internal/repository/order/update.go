package order

import (
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository/converter"
)

func (r *repository) Update(order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[order.OrderUUID]; !ok {
		return model.ErrOrderNotFound
	}
	r.orders[order.OrderUUID] = converter.OrderToModel(order)
	return nil
}
