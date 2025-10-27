package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/service"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

type api struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{orderService: orderService}
}

func (a *api) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return nil
}
