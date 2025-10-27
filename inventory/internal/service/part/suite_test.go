package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	repository *mocks.InventoryRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.repository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(
		s.repository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
