package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/vipshark78/microservices-course-homeworks/payment/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	service *mocks.PaymentService

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.service = mocks.NewPaymentService(s.T())

	s.api = NewApi(
		s.service,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
