package order

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *ServiceSuite) TestOrderPaySuccess() {
	var (
		orderUUID = uuid.New()
		order     = model.Order{
			OrderUUID:       orderUUID.String(),
			UserUUID:        orderUUID.String(),
			PartUuids:       []string{orderUUID.String()},
			TotalPrice:      300,
			TransactionUUID: orderUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
		payOrder = model.PayOrder{
			OrderUUID:     orderUUID.String(),
			UserUUID:      orderUUID.String(),
			PaymentMethod: model.PaymentMethodSBP,
		}
	)

	s.repository.On("Read", s.ctx, orderUUID.String()).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, payOrder).Return(orderUUID.String(), nil).Once()
	orderPaid := order
	orderPaid.Status = model.OrderStatusPAID
	orderPaid.TransactionUUID = orderUUID.String()
	s.repository.On("Update", orderPaid).Return(nil).Once()

	resp, err := s.service.OrderPay(s.ctx, orderUUID.String(), model.PaymentMethodSBP)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(orderUUID.String(), resp)
}

func (s *ServiceSuite) TestOrderPayErrorOrderNotFound() {
	orderUUID := uuid.New()

	s.repository.On("Read", s.ctx, orderUUID.String()).Return(model.Order{}, model.ErrOrderNotFound).Once()
	resp, err := s.service.OrderPay(s.ctx, orderUUID.String(), model.PaymentMethodSBP)

	s.NotNil(err)
	s.NotNil(resp)
	s.Equal(model.ErrOrderNotFound, err)
}

func (s *ServiceSuite) TestOrderPayErrorOrderPaid() {
	var (
		orderUUID = uuid.New()
		order     = model.Order{
			OrderUUID:       orderUUID.String(),
			UserUUID:        orderUUID.String(),
			PartUuids:       []string{orderUUID.String()},
			TotalPrice:      300,
			TransactionUUID: orderUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusPAID,
		}
	)

	s.repository.On("Read", s.ctx, orderUUID.String()).Return(order, nil).Once()
	resp, err := s.service.OrderPay(s.ctx, orderUUID.String(), model.PaymentMethodSBP)

	s.Require().Error(err)
	s.Require().NotNil(resp)
	s.Require().Equal(model.ErrOrderPaid, err)
}

func (s *ServiceSuite) TestOrderPayErrorOrderCancelled() {
	var (
		orderUUID = uuid.New()
		order     = model.Order{
			OrderUUID:       orderUUID.String(),
			UserUUID:        orderUUID.String(),
			PartUuids:       []string{orderUUID.String()},
			TotalPrice:      300,
			TransactionUUID: orderUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusCANCELLED,
		}
	)

	s.repository.On("Read", s.ctx, orderUUID.String()).Return(order, nil).Once()
	resp, err := s.service.OrderPay(s.ctx, orderUUID.String(), model.PaymentMethodSBP)

	s.Require().Error(err)
	s.Require().NotNil(resp)
	s.Require().Equal(model.ErrOrderCancelled, err)
}

func (s *ServiceSuite) TestOrderPayErrorOrderUpdateError() {
	var (
		orderUUID = uuid.New()
		order     = model.Order{
			OrderUUID:       orderUUID.String(),
			UserUUID:        orderUUID.String(),
			PartUuids:       []string{orderUUID.String()},
			TotalPrice:      300,
			TransactionUUID: orderUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
		payOrder = model.PayOrder{
			OrderUUID:     orderUUID.String(),
			UserUUID:      orderUUID.String(),
			PaymentMethod: model.PaymentMethodSBP,
		}
	)

	s.repository.On("Read", s.ctx, orderUUID.String()).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, payOrder).Return(orderUUID.String(), nil).Once()
	orderPaid := order
	orderPaid.Status = model.OrderStatusPAID
	orderPaid.TransactionUUID = orderUUID.String()
	s.repository.On("Update", orderPaid).Return(model.ErrOrderNotFound).Once()

	resp, err := s.service.OrderPay(s.ctx, orderUUID.String(), model.PaymentMethodSBP)

	s.Require().Error(err)
	s.Require().NotNil(resp)
	s.Require().Equal(model.ErrOrderNotFound, err)
}
