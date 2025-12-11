package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	authv1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, request *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if request == nil {
		return &authv1.LoginResponse{}, model.ErrInvalidCredentials
	}
	if request.Login == "" || request.Password == "" {
		return &authv1.LoginResponse{}, model.ErrInvalidCredentials
	}
	sessionUUID, err := a.service.Login(ctx, request.Login, request.Password)
	if err != nil {
		return &authv1.LoginResponse{}, err
	}

	return &authv1.LoginResponse{SessionUuid: sessionUUID}, nil
}
