package user

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/hash"
)

func (s *service) Register(ctx context.Context, userInfo model.UserRegistrationInfo) (string, error) {
	passHash, err := hash.Hash(userInfo.Password)
	if err != nil {
		return "", err
	}

	uuid, err := s.userRepo.Create(ctx, userInfo.Info, passHash)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
