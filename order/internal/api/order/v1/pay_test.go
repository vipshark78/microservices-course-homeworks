package v1

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestOrderPaySuccess() {
	var (
		orderUUID       = uuid.New()
		transactionUUID = uuid.NewString()
		payReq          = &order_v1.PayOrderRequest{
			PaymentMethod: order_v1.NewNilPaymentMethod(order_v1.PaymentMethodSBP),
		}
		payParams = order_v1.OrderPayParams{OrderUUID: orderUUID}
	)

	s.service.On("OrderPay", s.ctx, orderUUID.String(), model.PaymentMethodSBP).Return(transactionUUID, nil).Once()

	resp, err := s.api.OrderPay(s.ctx, payReq, payParams)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(transactionUUID, res.TransactionUUID.String())
}

func (s *APISuite) TestOrderPayNotFoundError() {
	var (
		orderUUID = uuid.New()
		payReq    = &order_v1.PayOrderRequest{
			PaymentMethod: order_v1.NewNilPaymentMethod(order_v1.PaymentMethodSBP),
		}
		payParams = order_v1.OrderPayParams{OrderUUID: orderUUID}
	)

	s.service.On("OrderPay", s.ctx, orderUUID.String(), model.PaymentMethodSBP).Return("", model.ErrOrderNotFound).Once()

	resp, err := s.api.OrderPay(s.ctx, payReq, payParams)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, res.Code)
}

func (s *APISuite) TestOrderPayBadRequestError() {
	resp, err := s.api.OrderPay(s.ctx, &order_v1.PayOrderRequest{}, order_v1.OrderPayParams{})

	s.Require().Error(err)
	badRequestErr, ok := resp.(*order_v1.BadRequestError)
	s.True(ok)
	s.Equal(400, badRequestErr.Code)
}
