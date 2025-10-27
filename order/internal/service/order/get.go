package order

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *service) OrderByUUID(ctx context.Context, orderUuid string) (model.Order, error) {
	return s.orderRepository.Read(ctx, orderUuid)
}
