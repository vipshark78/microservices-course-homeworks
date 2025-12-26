package auth

import (
	"context"
	"time"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/hash"
)

const ttl = time.Hour * 24

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(ctx, username)
	if err != nil {
		return "", err
	}

	err = hash.Compare(user.UserInfo.PasswordHash, password)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	session, err := model.NewSession(user.UserUUID, ttl)
	if err != nil {
		return "", err
	}
	sessionUuid, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		return "", err
	}
	return sessionUuid, nil
}
