package order

import (
	"sync"

	repomodel "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/model"
)

type repository struct {
	mu     sync.RWMutex
	orders map[string]repomodel.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]repomodel.Order),
	}
}
