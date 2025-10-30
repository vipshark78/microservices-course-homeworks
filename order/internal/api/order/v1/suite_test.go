package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	service *mocks.OrderService

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.service = mocks.NewOrderService(s.T())

	s.api = NewAPI(
		s.service,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
