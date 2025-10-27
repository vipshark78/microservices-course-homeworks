package order

import (
	"github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
