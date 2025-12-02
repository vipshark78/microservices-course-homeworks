package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	auth_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
	common_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/common/v1"
)

func ConvertServiceWhoamiRespToProtoModel(response model.WhoamiResponse) *auth_v1.WhoamiResponse {
	return &auth_v1.WhoamiResponse{
		Session: ConvertServiceSessionToProtoModel(response.Session),
		User:    ConvertServiceUserToProtoModel(response.User),
	}
}

func ConvertServiceSessionToProtoModel(session model.Session) *common_v1.Session {
	return &common_v1.Session{
		Uuid:      session.UUID,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
