package v1

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

// OrderCancel реализует метод OrderCancel из интерфейса API.
func (a *api) OrderCancel(ctx context.Context, params order_v1.OrderCancelParams) (order_v1.OrderCancelRes, error) {
	if err := uuid.Validate(params.OrderUUID.String()); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	err := a.orderService.OrderCancel(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderCancelled) {
			return &order_v1.ConflictError{Code: 409, Message: "Order has already been cancelled"}, nil
		}
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{Code: 404, Message: "Order Not Found"}, nil
		}
		if errors.Is(err, model.ErrOrderPaid) {
			return &order_v1.ConflictError{Code: 409, Message: "Order has already paid and cannot be cancelled"}, nil
		}
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}
	return &order_v1.OrderCancelNoContent{}, nil
}
