package order

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUuid string, partUuids []string) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{UUIDs: partUuids})
	if err != nil {
		return model.Order{}, err
	}

	price := s.calculatePrice(parts)

	return s.orderRepository.Insert(ctx, userUuid, partUuids, price)
}

func (s *service) calculatePrice(parts []model.Part) float64 {
	var price float64
	for _, part := range parts {
		price += part.Price
	}
	return price
}
