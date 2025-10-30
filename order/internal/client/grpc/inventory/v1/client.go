package v1

import inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"

type inventoryClient struct {
	client inventory_v1.InventoryServiceClient
}

func NewInventoryClient(client inventory_v1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{client: client}
}
