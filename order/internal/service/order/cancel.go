package order

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *service) OrderCancel(ctx context.Context, orderUuid string) error {
	order, err := s.orderRepository.Read(ctx, orderUuid)
	if err != nil {
		return err
	}

	switch order.Status {
	case model.OrderStatusCANCELLED:
		return model.ErrOrderCancelled
	case model.OrderStatusPAID:
		return model.ErrOrderPaid
	}
	order.Status = model.OrderStatusCANCELLED

	err = s.orderRepository.Update(order)
	if err != nil {
		return err
	}
	return nil
}
