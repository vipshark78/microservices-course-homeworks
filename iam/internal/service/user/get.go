package user

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, userUuid string) (model.User, error) {
	return s.userRepo.GetByUUID(ctx, userUuid)
}
