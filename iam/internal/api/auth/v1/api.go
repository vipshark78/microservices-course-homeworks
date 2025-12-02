package v1

import (
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/service"
	authV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAPI(service service.AuthService) *api {
	return &api{service: service}
}
