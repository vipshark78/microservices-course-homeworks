package v1

import (
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	client payment_v1.PaymentServiceClient
}

func NewPaymentClient(client payment_v1.PaymentServiceClient) *paymentClient {
	return &paymentClient{client: client}
}
