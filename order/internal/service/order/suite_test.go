package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc/mocks"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	repository *mocks.OrderRepository

	inventoryClient *grpcMocks.InventoryClient
	paymentClient   *grpcMocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.repository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())

	s.service = NewService(
		s.repository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
