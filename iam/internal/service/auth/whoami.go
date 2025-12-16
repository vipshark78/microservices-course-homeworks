package auth

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUuid string) (model.WhoamiResponse, error) {
	session, err := s.sessionRepo.Get(ctx, sessionUuid)
	if err != nil {
		return model.WhoamiResponse{}, err
	}

	user, err := s.userRepo.GetByUUID(ctx, session.UUID())
	if err != nil {
		return model.WhoamiResponse{}, err
	}

	return model.WhoamiResponse{
		Session: session,
		User:    user,
	}, nil
}
