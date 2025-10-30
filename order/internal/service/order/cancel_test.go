package order

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *ServiceSuite) TestOrderCancelSuccess() {
	var (
		reqUUID = uuid.New()
		order   = model.Order{
			OrderUUID:       reqUUID.String(),
			UserUUID:        reqUUID.String(),
			PartUuids:       []string{reqUUID.String()},
			TotalPrice:      300,
			TransactionUUID: reqUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.repository.On("Read", s.ctx, reqUUID.String()).Return(order, nil).Once()
	orderCancelled := order
	orderCancelled.Status = model.OrderStatusCANCELLED
	s.repository.On("Update", orderCancelled).Return(nil).Once()

	err := s.service.OrderCancel(s.ctx, reqUUID.String())

	s.NoError(err)
}

func (s *ServiceSuite) TestOrderCancelNotFound() {
	reqUUID := uuid.New()

	s.repository.On("Read", s.ctx, reqUUID.String()).Return(model.Order{}, model.ErrOrderNotFound).Once()

	err := s.service.OrderCancel(s.ctx, reqUUID.String())

	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderNotFound, err)
}

func (s *ServiceSuite) TestOrderCancelErrorOrderCancelled() {
	var (
		reqUUID = uuid.New()
		order   = model.Order{
			OrderUUID:       reqUUID.String(),
			UserUUID:        reqUUID.String(),
			PartUuids:       []string{reqUUID.String()},
			TotalPrice:      300,
			TransactionUUID: reqUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusCANCELLED,
		}
	)

	s.repository.On("Read", s.ctx, reqUUID.String()).Return(order, nil).Once()

	err := s.service.OrderCancel(s.ctx, reqUUID.String())

	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderCancelled, err)
}

func (s *ServiceSuite) TestOrderCancelErrorOrderPaid() {
	var (
		reqUUID = uuid.New()
		order   = model.Order{
			OrderUUID:       reqUUID.String(),
			UserUUID:        reqUUID.String(),
			PartUuids:       []string{reqUUID.String()},
			TotalPrice:      300,
			TransactionUUID: reqUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusPAID,
		}
	)

	s.repository.On("Read", s.ctx, reqUUID.String()).Return(order, nil).Once()

	err := s.service.OrderCancel(s.ctx, reqUUID.String())

	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderPaid, err)
}

func (s *ServiceSuite) TestOrderCancelUpdateError() {
	var (
		reqUUID = uuid.New()
		order   = model.Order{
			OrderUUID:       reqUUID.String(),
			UserUUID:        reqUUID.String(),
			PartUuids:       []string{reqUUID.String()},
			TotalPrice:      300,
			TransactionUUID: reqUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.repository.On("Read", s.ctx, reqUUID.String()).Return(order, nil).Once()

	orderCancelled := order
	orderCancelled.Status = model.OrderStatusCANCELLED

	s.repository.On("Update", orderCancelled).Return(model.ErrOrderNotFound).Once()

	err := s.service.OrderCancel(s.ctx, reqUUID.String())

	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderNotFound, err)
}
