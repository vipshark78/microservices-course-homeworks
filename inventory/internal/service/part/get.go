package part

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
)

// GetPart получение детали по UUID.
func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	return s.inventoryRepository.GetPart(ctx, uuid)
}
