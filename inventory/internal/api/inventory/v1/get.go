package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

// GetPart возвращает информацию о детали по UUID
func (a *api) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	if req == nil {
		return nil, model.ErrInvalidRequest
	}
	if req.Uuid == "" {
		return nil, model.ErrInvalidUUID
	}
	uuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, model.ErrInvalidUUID
	}

	part, err := a.inventoryService.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &inventory_v1.GetPartResponse{
		Part: converter.ModelToPart(part),
	}, nil
}
