package service

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payOrder model.PayOrder) (string, error)
}
