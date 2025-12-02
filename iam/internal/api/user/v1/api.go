package v1

import (
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/service"
	userV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer
	service service.UserService
}

func NewAPI(service service.UserService) *api {
	return &api{service: service}
}
