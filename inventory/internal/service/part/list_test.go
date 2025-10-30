package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	var (
		uuids  = []string{gofakeit.UUID(), gofakeit.UUID()}
		filter = model.PartsFilter{UUIDs: uuids}
		parts  = []model.Part{
			{UUID: uuids[0]},
			{UUID: uuids[1]},
		}
	)

	s.repository.On("ListParts", s.ctx, filter).Return(parts, nil)
	resp, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)
	s.Equal(parts, resp)
}

func (s *ServiceSuite) TestListPartsNotFoundError() {
	var (
		uuids  = []string{gofakeit.UUID(), gofakeit.UUID()}
		filter = model.PartsFilter{UUIDs: uuids}
	)

	s.repository.On("ListParts", s.ctx, filter).Return(nil, model.ErrPartNotFound)
	resp, err := s.service.ListParts(s.ctx, filter)

	s.Error(err)
	s.ErrorIs(err, model.ErrPartNotFound)
	s.Nil(resp)
}
