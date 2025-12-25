package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc/mocks"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/metrics"
	repoMocks "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/mocks"
	serviceMocks "github.com/vipshark78/microservices-course-homeworks/order/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	repository *repoMocks.OrderRepository

	inventoryClient *grpcMocks.InventoryClient
	paymentClient   *grpcMocks.PaymentClient
	producerService *serviceMocks.OrderProducerService

	service *service
}

func (s *ServiceSuite) SetupTest() {
	_ = metrics.InitMetrics()
	s.ctx = context.Background()

	s.repository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())
	s.producerService = serviceMocks.NewOrderProducerService(s.T())
	s.service = NewService(
		s.repository,
		s.inventoryClient,
		s.paymentClient,
		s.producerService,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
