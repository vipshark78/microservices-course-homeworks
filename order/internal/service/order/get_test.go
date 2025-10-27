package order

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *ServiceSuite) TestOrderByUUIDSuccess() {
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

	s.repository.On("Read", s.ctx, reqUUID.String()).Return(order, nil)
	resp, err := s.service.OrderByUUID(s.ctx, reqUUID.String())

	s.NoError(err)
	s.Equal(order, resp)
}

func (s *ServiceSuite) TestOrderByUUIDNotFound() {
	reqUUID := uuid.New()

	s.repository.On("Read", s.ctx, reqUUID.String()).Return(model.Order{}, model.ErrOrderNotFound)

	resp, err := s.service.OrderByUUID(s.ctx, reqUUID.String())
	s.Require().NotNil(err)
	s.Require().NotNil(resp)
	s.Require().Equal(model.ErrOrderNotFound, err)
}
