package auth

import "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository"

type service struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}
