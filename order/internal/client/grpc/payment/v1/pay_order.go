package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/client/converter"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (p *paymentClient) PayOrder(ctx context.Context, payOrder model.PayOrder) (string, error) {
	resp, err := p.client.PayOrder(ctx, converter.ModelPayOrderToProtoPayOrder(payOrder))
	if err != nil {
		return "", err
	}
	return resp.TransactionUuid, nil
}
