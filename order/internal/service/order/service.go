package order

import (
	"github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository"
	orderService "github.com/vipshark78/microservices-course-homeworks/order/internal/service"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
	producer        orderService.OrderProducerService
}

func NewService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, producer orderService.OrderProducerService) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		producer:        producer,
	}
}
