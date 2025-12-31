package v1

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/converter"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	authV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	if req.SessionUuid == "" {
		return &authV1.WhoamiResponse{}, model.ErrInvalidSessionUUID
	}
	resp, err := a.service.Whoami(ctx, req.SessionUuid)
	if err != nil {
		return &authV1.WhoamiResponse{}, err
	}

	return converter.ConvertServiceWhoamiRespToProtoModel(resp), nil
}
