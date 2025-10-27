package v1

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/payment/internal/model"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		orderUUID = uuid.New()
		payOrder  = model.PayOrder{
			OrderUUID:     orderUUID.String(),
			UserUUID:      orderUUID.String(),
			PaymentMethod: model.PaymentMethodCREDITCARD,
		}
		payOrderRequest = &payment_v1.PayOrderRequest{
			OrderUuid:     orderUUID.String(),
			UserUuid:      orderUUID.String(),
			PaymentMethod: payment_v1.PaymentMethod_CREDIT_CARD,
		}
	)

	s.service.On("PayOrder", s.ctx, payOrder).Return(orderUUID.String(), nil).Once()

	resp, err := s.api.PayOrder(s.ctx, payOrderRequest)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(orderUUID.String(), resp.TransactionUuid)
}
