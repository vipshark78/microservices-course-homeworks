package repository

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
)

type UserRepository interface {
	GetByUUID(ctx context.Context, userUUID string) (model.User, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
	Create(ctx context.Context, info model.UserInfo, hash string) (string, error)
}

type SessionRepository interface {
	Get(ctx context.Context, sessionUUID string) (model.Session, error)
	Create(ctx context.Context, userUUID string) (string, error)
	AddToUserSet(ctx context.Context, userUUID, sessionUUID string) error
}
