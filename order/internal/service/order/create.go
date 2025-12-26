package order

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	orderMetrics "github.com/vipshark78/microservices-course-homeworks/order/internal/metrics"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUuid string, partUuids []string) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{UUIDs: partUuids})
	if err != nil {
		return model.Order{}, err
	}

	price := s.calculatePrice(parts)

	orderMetrics.OrdersTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("status", "success")),
	)
	return s.orderRepository.Insert(ctx, userUuid, partUuids, price, model.OrderStatusPENDINGPAYMENT)
}

func (s *service) calculatePrice(parts []model.Part) float64 {
	var price float64
	for _, part := range parts {
		price += part.Price
	}
	return price
}
