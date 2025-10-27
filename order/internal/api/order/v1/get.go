package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/converter"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

// OrderByUUID реализует получение заказа по UUID.
func (a *api) OrderByUUID(ctx context.Context, params order_v1.OrderByUUIDParams) (order_v1.OrderByUUIDRes, error) {
	if err := uuid.Validate(params.OrderUUID.String()); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	order, err := a.orderService.OrderByUUID(ctx, params.OrderUUID.String())
	if err != nil {
		return &order_v1.NotFoundError{Code: 404, Message: "Order Not Found"}, nil
	}

	return &order_v1.GetOrderResponse{AllOf: order_v1.NewOptOrderDto(converter.ConvertModelToOrder(order))}, nil
}
