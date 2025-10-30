package v1

import (
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/service"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
