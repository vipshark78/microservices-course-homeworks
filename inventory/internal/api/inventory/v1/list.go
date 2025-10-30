package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

// ListParts возвращает список деталей с возможностью фильтрации
func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	if req == nil || req.Filter == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}
	modelParts, err := a.inventoryService.ListParts(ctx, converter.PartsFilterToModel(req.Filter))
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "no parts")
		}
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}
	parts := converter.ModelsToParts(modelParts)
	return &inventory_v1.ListPartsResponse{
		Parts: parts,
	}, nil
}
