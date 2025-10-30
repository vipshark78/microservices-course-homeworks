package repository

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, string string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
