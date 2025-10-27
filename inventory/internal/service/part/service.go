package part

import "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository"

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
