package v1

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestOrderByUUIDSuccess() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		totalPrice      = float64(20.5)
		status          = model.OrderStatusPENDINGPAYMENT
		paymentMethod   = model.PaymentMethodSBP
		transactionUUID = uuid.New()
		partUUIDs       = []string{uuid.New().String(), uuid.New().String()}
		reqParams       = order_v1.OrderByUUIDParams{
			OrderUUID: orderUUID,
		}
		respOrder = model.Order{
			OrderUUID:       orderUUID.String(),
			UserUUID:        userUUID.String(),
			PartUuids:       partUUIDs,
			TotalPrice:      totalPrice,
			TransactionUUID: transactionUUID.String(),
			PaymentMethod:   paymentMethod,
			Status:          status,
		}
	)

	s.service.On("OrderByUUID", s.ctx, orderUUID.String()).Return(respOrder, nil).Once()

	resp, err := s.api.OrderByUUID(s.ctx, reqParams)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.GetOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(order_v1.NewOptOrderDto(converter.ConvertModelToOrder(respOrder)), res.AllOf)
}

func (s *APISuite) TestOrderByUUIDNotFound() {
	var (
		orderUUID = uuid.New()
		reqParams = order_v1.OrderByUUIDParams{
			OrderUUID: orderUUID,
		}
	)

	s.service.On("OrderByUUID", s.ctx, orderUUID.String()).Return(model.Order{}, model.ErrOrderNotFound).Once()

	resp, err := s.api.OrderByUUID(s.ctx, reqParams)
	s.Require().Nil(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, res.Code)
}
