package app

import (
	"context"

	paymentV1API "github.com/vipshark78/microservices-course-homeworks/payment/internal/api/payment/v1"
	"github.com/vipshark78/microservices-course-homeworks/payment/internal/service"
	paymentService "github.com/vipshark78/microservices-course-homeworks/payment/internal/service/payment"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API   payment_v1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) payment_v1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewApi(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewPaymentService()
	}

	return d.paymentService
}
