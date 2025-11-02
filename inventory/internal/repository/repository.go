package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	repoModel "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
	InsertPart(ctx context.Context, part repoModel.Part) error
}
