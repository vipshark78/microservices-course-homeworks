package payment

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/payment/internal/model"
)

func (s *ServiceSuite) TestOrderPaySuccess() {
	var (
		orderUUID = uuid.New()
		payOrder  = model.PayOrder{
			OrderUUID:     orderUUID.String(),
			UserUUID:      orderUUID.String(),
			PaymentMethod: model.PaymentMethodSBP,
		}
	)

	res, err := s.service.PayOrder(s.ctx, payOrder)

	s.NoError(err)
	s.Require().NotNil(res)
}
