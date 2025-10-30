package service

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
