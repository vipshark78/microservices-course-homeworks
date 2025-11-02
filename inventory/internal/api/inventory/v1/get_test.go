package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	business_errors "github.com/vipshark78/microservices-course-homeworks/shared/pkg/errors"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	var (
		reqUUID = uuid.New()
		req     = &inventory_v1.GetPartRequest{Uuid: reqUUID.String()}
		part    = model.Part{
			UUID:          req.Uuid,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Word(),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      model.Category(gofakeit.Word()),
			Dimensions: &model.Dimensions{
				Length: gofakeit.Float64Range(10, 100),
				Width:  gofakeit.Float64Range(10, 100),
				Height: gofakeit.Float64Range(10, 100),
				Weight: gofakeit.Float64Range(1, 5),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{gofakeit.EmojiTag(), gofakeit.EmojiTag(), gofakeit.EmojiTag()},
			Metadata: map[string]model.Value{
				gofakeit.Word(): {StringValue: lo.ToPtr(gofakeit.Word())},
				gofakeit.Word(): {Int64Value: lo.ToPtr(int64(gofakeit.Number(1, 9)))},
				gofakeit.Word(): {BooleanValue: lo.ToPtr(gofakeit.Bool())},
				gofakeit.Word(): {DoubleValue: lo.ToPtr(gofakeit.Float64Range(1, 10))},
			},
			CreatedAt: lo.ToPtr(time.Now()),
			UpdatedAt: lo.ToPtr(time.Now()),
		}
		protoPart = converter.ModelToPart(part)
	)

	s.service.On("GetPart", s.ctx, reqUUID).Return(part, nil).Once()

	resp, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(protoPart, resp.Part)
}

func (s *APISuite) TestGetPartBadRequestErrorEmptyUUID() {
	resp, err := s.api.GetPart(s.ctx, &inventory_v1.GetPartRequest{})

	s.Require().Error(err)
	s.Nil(resp)
	businessErr := business_errors.GetBusinessError(err)
	s.Equal(business_errors.InvalidArgumentErrCode, businessErr.Code())
}

func (s *APISuite) TestGetPartBadRequestErrorEmptyReq() {
	resp, err := s.api.GetPart(s.ctx, nil)

	s.Require().Error(err)
	s.Nil(resp)
	businessErr := business_errors.GetBusinessError(err)
	s.Equal(business_errors.InvalidArgumentErrCode, businessErr.Code())
}

func (s *APISuite) TestGetPartNotFoundError() {
	reqUUID := uuid.New()
	req := &inventory_v1.GetPartRequest{Uuid: reqUUID.String()}

	s.service.On("GetPart", s.ctx, reqUUID).Return(model.Part{}, model.ErrPartNotFound).Once()

	resp, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Nil(resp)
	businessErr := business_errors.GetBusinessError(err)
	s.Equal(business_errors.NotFoundErrCode, businessErr.Code())
}

func (s *APISuite) TestGetPartInternalServerError() {
	reqUUID := uuid.New()
	req := &inventory_v1.GetPartRequest{Uuid: reqUUID.String()}
	sErr := gofakeit.Error()

	s.service.On("GetPart", s.ctx, reqUUID).Return(model.Part{}, sErr).Once()

	resp, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Nil(resp)
	businessErr := business_errors.GetBusinessError(err)
	s.Nil(businessErr)
}
