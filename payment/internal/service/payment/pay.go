package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/payment/internal/model"
)

func (p *paymentService) PayOrder(ctx context.Context, payOrder model.PayOrder) (string, error) {
	transactionUUID := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", transactionUUID)
	return transactionUUID, nil
}
