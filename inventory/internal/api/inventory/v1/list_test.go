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

func (s *APISuite) TestListPartsSuccess() {
	var (
		names     = []string{gofakeit.Name(), gofakeit.Name()}
		reqFilter = &inventory_v1.PartsFilter{
			Names: names,
		}
		req    = &inventory_v1.ListPartsRequest{Filter: reqFilter}
		filter = converter.PartsFilterToModel(reqFilter)
		parts  = []model.Part{
			{
				UUID:          gofakeit.UUID(),
				Name:          names[0],
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
			},
			{
				UUID:          gofakeit.UUID(),
				Name:          names[1],
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
			},
		}
		protoParts = converter.ModelsToParts(parts)
	)

	s.service.On("ListParts", s.ctx, filter).Return(parts, nil).Once()

	resp, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(protoParts, resp.Parts)
}

func (s *APISuite) TestListPartBadRequestError() {
	resp, err := s.api.ListParts(s.ctx, &inventory_v1.ListPartsRequest{})

	s.Require().Error(err)
	s.Nil(resp)
	s.Equal(codes.InvalidArgument, status.Code(err))
}

func (s *APISuite) TestListPartsNotFoundError() {
	req := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: []string{gofakeit.UUID()},
		},
	}

	filter := converter.PartsFilterToModel(req.Filter)

	s.service.On("ListParts", s.ctx, filter).Return(nil, model.ErrPartNotFound).Once()
	resp, err := s.api.ListParts(s.ctx, req)
	s.Require().Error(err)
	s.Nil(resp)
	s.Equal(codes.NotFound, status.Code(err))
}

func (s *APISuite) TestListPartsInternalError() {
	req := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			ManufacturerCountries: []string{gofakeit.Country()},
		},
	}

	filter := converter.PartsFilterToModel(req.Filter)

	sErr := gofakeit.Error()

	s.service.On("ListParts", s.ctx, filter).Return(nil, sErr).Once()

	resp, err := s.api.ListParts(s.ctx, req)

	s.Require().Error(err)
	s.Nil(resp)
	s.Equal(codes.Internal, status.Code(err))
}
