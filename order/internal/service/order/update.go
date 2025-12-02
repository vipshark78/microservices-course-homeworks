package order

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *service) OrderUpdateStatus(ctx context.Context, orderUuid string, orderStatus model.OrderStatus) error {
	order, err := s.orderRepository.Read(ctx, orderUuid)
	if err != nil {
		return err
	}

	order.Status = orderStatus

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}
	return nil
}
