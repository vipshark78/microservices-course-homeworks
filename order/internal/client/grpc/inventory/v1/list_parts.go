package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/client/converter"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

func (i *inventoryClient) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	resp, err := i.client.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: converter.ModelToPartsFilter(filter),
	})
	if err != nil {
		return nil, err
	}
	return converter.ProtoPartsToModelParts(resp.Parts), nil
}
