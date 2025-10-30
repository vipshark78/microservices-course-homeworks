package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/payment/internal/model"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.PayOrder(ctx, model.PayOrder{OrderUUID: req.OrderUuid, UserUUID: req.UserUuid, PaymentMethod: req.PaymentMethod.String()})
	if err != nil {
		return nil, err
	}
	return &payment_v1.PayOrderResponse{TransactionUuid: transactionUUID}, err
}
