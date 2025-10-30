package order

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *service) OrderPay(ctx context.Context, orderUuid, paymentMethod string) (string, error) {
	order, err := s.orderRepository.Read(ctx, orderUuid)
	if err != nil {
		return "", err
	}

	switch order.Status {
	case model.OrderStatusCANCELLED:
		return "", model.ErrOrderCancelled
	case model.OrderStatusPAID:
		return "", model.ErrOrderPaid
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, model.PayOrder{OrderUUID: orderUuid, UserUUID: order.UserUUID, PaymentMethod: paymentMethod})
	if err != nil {
		return "", err
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = paymentMethod
	order.TransactionUUID = transactionUUID

	err = s.orderRepository.Update(order)
	if err != nil {
		return "", err
	}
	return transactionUUID, nil
}
