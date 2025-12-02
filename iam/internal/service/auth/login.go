package auth

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/hash"
)

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(ctx, username)
	if err != nil {
		return "", err
	}

	err = hash.Compare(user.UserInfo.PasswordHash, password)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	sessionUuid, err := s.sessionRepo.Create(ctx, user.UserUUID)
	if err != nil {
		return "", err
	}
	return sessionUuid, nil
}
