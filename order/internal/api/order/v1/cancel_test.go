package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestOrderCancelSuccess() {
	var (
		orderUUID    = uuid.New()
		cancelParams = order_v1.OrderCancelParams{OrderUUID: orderUUID}
	)

	s.service.On("OrderCancel", s.ctx, orderUUID.String()).Return(nil).Once()

	resp, err := s.api.OrderCancel(s.ctx, cancelParams)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(&order_v1.OrderCancelNoContent{}, resp)
}

func (s *APISuite) TestOrderCancelNotFoundError() {
	var (
		orderUUID    = uuid.New()
		cancelParams = order_v1.OrderCancelParams{OrderUUID: orderUUID}
	)

	s.service.On("OrderCancel", s.ctx, orderUUID.String()).Return(model.ErrOrderNotFound).Once()

	resp, err := s.api.OrderCancel(s.ctx, cancelParams)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, res.Code)
}

func (s *APISuite) TestOrderCancelConflictErrorOrderCancelled() {
	var (
		orderUUID    = uuid.New()
		cancelParams = order_v1.OrderCancelParams{OrderUUID: orderUUID}
	)

	s.service.On("OrderCancel", s.ctx, orderUUID.String()).Return(model.ErrOrderCancelled).Once()

	resp, err := s.api.OrderCancel(s.ctx, cancelParams)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(409, res.Code)
}

func (s *APISuite) TestOrderCancelConflictErrorOrderPaid() {
	var (
		orderUUID    = uuid.New()
		cancelParams = order_v1.OrderCancelParams{OrderUUID: orderUUID}
	)

	s.service.On("OrderCancel", s.ctx, orderUUID.String()).Return(model.ErrOrderPaid).Once()

	resp, err := s.api.OrderCancel(s.ctx, cancelParams)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(409, res.Code)
}

func (s *APISuite) TestOrderCancelInternalServerError() {
	var (
		orderUUID    = uuid.New()
		cancelParams = order_v1.OrderCancelParams{OrderUUID: orderUUID}
	)

	s.service.On("OrderCancel", s.ctx, orderUUID.String()).Return(gofakeit.Error()).Once()

	resp, err := s.api.OrderCancel(s.ctx, cancelParams)
	s.Require().Error(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(500, res.Code)
}
