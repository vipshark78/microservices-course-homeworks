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

// GetPart возвращает информацию о детали по UUID
func (a *api) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request: uuid is empty")
	}
	part, err := a.inventoryService.GetPart(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "no part with uuid '%s'", req.Uuid)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inventory_v1.GetPartResponse{
		Part: converter.ModelToPart(part),
	}, nil
}
