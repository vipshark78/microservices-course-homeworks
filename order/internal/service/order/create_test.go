package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
)

func (s *ServiceSuite) TestOrderCreateSuccess() {
	var (
		reqUUID          = uuid.New()
		partUUID1        = uuid.New()
		partUUID2        = uuid.New()
		partsUUIDs       = model.PartsFilter{UUIDs: []string{partUUID1.String(), partUUID2.String()}}
		partsUUIDsString = []string{partUUID1.String(), partUUID2.String()}
		price1           = float64(500)
		price2           = float64(500)
		parts            = []model.Part{
			{UUID: partUUID1.String(), Price: price1},
			{UUID: partUUID2.String(), Price: price2},
		}
		totalPrice = price1 + price2
		order      = model.Order{
			OrderUUID:       reqUUID.String(),
			UserUUID:        reqUUID.String(),
			PartUuids:       partsUUIDsString,
			TotalPrice:      totalPrice,
			TransactionUUID: reqUUID.String(),
			PaymentMethod:   model.PaymentMethodSBP,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, partsUUIDs).Return(parts, nil).Once()
	s.repository.On("Insert", s.ctx, reqUUID.String(), partsUUIDsString, totalPrice).Return(order, nil).Once()

	res, err := s.service.CreateOrder(s.ctx, reqUUID.String(), partsUUIDsString)

	s.NoError(err)
	s.Equal(order, res)
}

func (s *ServiceSuite) TestOrderCreateErrorPartsNotFound() {
	var (
		reqUUID          = uuid.New()
		partUUID1        = uuid.New()
		partUUID2        = uuid.New()
		partsUUIDs       = model.PartsFilter{UUIDs: []string{partUUID1.String(), partUUID2.String()}}
		partsUUIDsString = []string{partUUID1.String(), partUUID2.String()}
	)

	s.inventoryClient.On("ListParts", s.ctx, partsUUIDs).Return(nil, gofakeit.Error()).Once()
	// s.repository.On("Insert", s.ctx, reqUUID.String(), partsUUIDsString, totalPrice).Return(order, nil).Once()

	res, err := s.service.CreateOrder(s.ctx, reqUUID.String(), partsUUIDsString)

	s.NotNil(err)
	s.NotNil(res)
	s.Equal(model.Order{}, res)
}
