package v1

import (
	"github.com/vipshark78/microservices-course-homeworks/payment/internal/service"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewApi(service service.PaymentService) *api {
	return &api{paymentService: service}
}
