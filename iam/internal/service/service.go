package service

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Whoami(ctx context.Context, sessionUuid string) (model.WhoamiResponse, error)
}

type UserService interface {
	Register(ctx context.Context, userInfo model.UserRegistrationInfo) (string, error)
	GetUser(ctx context.Context, userUuid string) (model.User, error)
}
