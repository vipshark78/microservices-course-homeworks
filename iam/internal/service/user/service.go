package user

import "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository"

type service struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *service {
	return &service{
		userRepo: userRepo,
	}
}
