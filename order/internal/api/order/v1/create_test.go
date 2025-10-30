package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestOrderCreateBadRequestError() {
	resp, err := s.api.CreateOrder(s.ctx, &order_v1.CreateOrderRequest{})

	s.Require().Error(err)
	s.NotNil(resp)
	res, ok := resp.(*order_v1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(400, res.Code)
}

func (s *APISuite) TestOrderCreateSuccess() {
	var (
		userUUID  = uuid.New()
		orderUUID = uuid.New()
		partUUIDs = []uuid.UUID{uuid.New(), uuid.New()}
		req       = &order_v1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUUIDs,
		}
		partUUIDsStrings = []string{partUUIDs[0].String(), partUUIDs[1].String()}
		respOrder        = model.Order{OrderUUID: orderUUID.String(), UserUUID: userUUID.String(), PartUuids: partUUIDsStrings}
	)

	s.service.On("CreateOrder", s.ctx, userUUID.String(), partUUIDsStrings).Return(respOrder, nil).Once()

	resp, err := s.api.CreateOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.CreateOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(uuid.MustParse(respOrder.OrderUUID), res.OrderUUID)
}

func (s *APISuite) TestOrderCreateInternalServerError() {
	var (
		userUUID  = uuid.New()
		partUUIDs = []uuid.UUID{uuid.New(), uuid.New()}
		req       = &order_v1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUUIDs,
		}
		partUUIDsStrings = []string{partUUIDs[0].String(), partUUIDs[1].String()}
	)

	s.service.On("CreateOrder", s.ctx, userUUID.String(), partUUIDsStrings).Return(model.Order{}, gofakeit.Error()).Once()

	resp, err := s.api.CreateOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().NotNil(resp)
	res, ok := resp.(*order_v1.InternalServerError)
	s.Require().True(ok)
	s.Require().Equal(500, res.Code)
}
