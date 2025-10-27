package part

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		reqUUID = uuid.New()
		order   = model.Part{UUID: reqUUID.String()}
	)

	s.repository.On("GetPart", s.ctx, reqUUID.String()).Return(order, nil)
	resp, err := s.service.GetPart(s.ctx, reqUUID.String())

	s.NoError(err)
	s.Equal(order, resp)
}

func (s *ServiceSuite) TestGetPartNotFoundError() {
	reqUUID := uuid.New()

	s.repository.On("GetPart", s.ctx, reqUUID.String()).Return(model.Part{}, model.ErrPartNotFound)
	resp, err := s.service.GetPart(s.ctx, reqUUID.String())

	s.Error(err)
	s.ErrorIs(err, model.ErrPartNotFound)
	s.Equal(model.Part{}, resp)
}
