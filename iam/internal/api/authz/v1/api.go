package v1

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/service"
)

type api struct {
	authv3.UnimplementedAuthorizationServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
