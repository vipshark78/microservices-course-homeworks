package order

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	orderMetrics "github.com/vipshark78/microservices-course-homeworks/order/internal/metrics"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/tracing"
)

func (s *service) OrderPay(ctx context.Context, orderUuid string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Read(ctx, orderUuid)
	if err != nil {
		return "", err
	}

	switch order.Status {
	case model.OrderStatusCANCELLED:
		return "", model.ErrOrderCancelled
	case model.OrderStatusPAID:
		return "", model.ErrOrderPaid
	}

	ctx, span := tracing.StartSpan(ctx, "order.call_payment",
		trace.WithAttributes(
			attribute.String("order.uuid", order.OrderUUID),
		),
	)
	defer span.End()

	transactionUUID, err := s.paymentClient.PayOrder(ctx, model.PayOrder{OrderUUID: orderUuid, UserUUID: order.UserUUID, PaymentMethod: paymentMethod})
	if err != nil {
		return "", err
	}

	span.SetAttributes(
		attribute.String("payment.transactionUUID", transactionUUID),
		attribute.String("payment.status", string(order.Status)),
		attribute.Float64("payment.totalPrice", order.TotalPrice),
	)

	err = s.producer.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.New().String(),
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   string(paymentMethod),
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		return "", err
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = &paymentMethod
	order.TransactionUUID = &transactionUUID

	orderMetrics.OrdersRevenueTotal.Add(ctx, order.TotalPrice,
		metric.WithAttributes(
			attribute.String("currency", "RUB"),
			attribute.String("payment_method", string(*order.PaymentMethod)),
		),
	)

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return "", err
	}
	return transactionUUID, nil
}
