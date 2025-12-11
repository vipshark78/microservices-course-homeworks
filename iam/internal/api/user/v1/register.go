package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	userV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	if req.Info == nil || req.Info.Info == nil {
		return &userV1.RegisterResponse{}, model.ErrInvalidUserInfo
	}
	if req.Info.Password == "" || req.Info.Info.Login == "" || req.Info.Info.Email == "" {
		return &userV1.RegisterResponse{}, model.ErrInvalidRequest
	}

	userUUID, err := a.service.Register(ctx, converter.ConvertProtoUserRegInfoToServiceModel(req.Info))
	if err != nil {
		return &userV1.RegisterResponse{}, err
	}

	return &userV1.RegisterResponse{UserUuid: userUUID}, nil
}
