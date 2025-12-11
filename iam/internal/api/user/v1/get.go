package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	userV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/user/v1"
)

func (a *api) Get(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	if req.UserUuid == "" {
		return nil, model.ErrInvalidUserUUID
	}

	user, err := a.service.GetUser(ctx, req.UserUuid)
	if err != nil {
		return nil, err
	}

	return &userV1.GetUserResponse{
		User: converter.ConvertServiceUserToProtoModel(user),
	}, nil
}
