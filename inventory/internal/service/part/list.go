package part

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
)

// ListParts возвращает список деталей по фильтру.
func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	return s.inventoryRepository.ListParts(ctx, filter)
}
