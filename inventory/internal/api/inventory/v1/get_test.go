package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	var (
		req  = &inventory_v1.GetPartRequest{Uuid: gofakeit.UUID()}
		part = model.Part{
			UUID:          req.Uuid,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Word(),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      gofakeit.Word(),
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
			Metadata: map[string]*model.Value{
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

	s.service.On("GetPart", s.ctx, req.Uuid).Return(part, nil).Once()

	resp, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(protoPart, resp.Part)
}

func (s *APISuite) TestGetPartBadRequestErrorEmptyUUID() {
	resp, err := s.api.GetPart(s.ctx, &inventory_v1.GetPartRequest{})

	s.Require().Error(err)
	s.Nil(resp)
	s.Equal(codes.InvalidArgument, status.Code(err))
}

func (s *APISuite) TestGetPartBadRequestErrorEmptyReq() {
	resp, err := s.api.GetPart(s.ctx, nil)

	s.Require().Error(err)
	s.Nil(resp)
	s.Equal(codes.InvalidArgument, status.Code(err))
}

func (s *APISuite) TestGetPartNotFoundError() {
	req := &inventory_v1.GetPartRequest{Uuid: gofakeit.UUID()}

	s.service.On("GetPart", s.ctx, req.Uuid).Return(model.Part{}, model.ErrPartNotFound).Once()

	resp, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Nil(resp)
	s.Equal(codes.NotFound, status.Code(err))
	s.Contains(err.Error(), req.Uuid)
}

func (s *APISuite) TestGetPartInternalServerError() {
	req := &inventory_v1.GetPartRequest{Uuid: gofakeit.UUID()}
	sErr := gofakeit.Error()

	s.service.On("GetPart", s.ctx, req.Uuid).Return(model.Part{}, sErr).Once()

	resp, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Nil(resp)
	s.Equal(codes.Internal, status.Code(err))
}
