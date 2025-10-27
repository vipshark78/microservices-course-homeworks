package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/converter"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

// OrderPay реализует обработку запроса на оплату заказа.
func (a *api) OrderPay(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.OrderPayParams) (order_v1.OrderPayRes, error) {
	if err := a.validateOrderPayRequest(req, params); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	id, err := a.orderService.OrderPay(ctx, params.OrderUUID.String(), converter.ConvertPaymentMethodToModel(req.PaymentMethod))
	if err != nil {
		return &order_v1.NotFoundError{Code: 404, Message: "Order Not Found"}, nil
	}
	return &order_v1.PayOrderResponse{TransactionUUID: uuid.MustParse(id)}, nil
}

// validateOrderPayRequest выполняет валидацию входящих данных.
func (a *api) validateOrderPayRequest(req *order_v1.PayOrderRequest, params order_v1.OrderPayParams) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации запроса: %w", err)
	}

	if err := uuid.Validate(params.OrderUUID.String()); err != nil {
		return fmt.Errorf("ошибка валидации UUID заказа: %w", err)
	}

	if req.PaymentMethod.IsNull() {
		return fmt.Errorf("не указан метод оплаты")
	}
	return nil
}
