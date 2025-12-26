package converter

import (
	"time"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	repoModel "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/model"
)

func ConvertRepoSessionModelToModel(repoSession repoModel.Session) model.Session {
	return model.NewSessionFromValues(
		repoSession.UserUUID,
		time.Unix(repoSession.CreatedAt, 0),
		time.Unix(repoSession.UpdatedAt, 0),
		time.Unix(repoSession.ExpiresAt, 0),
	)
}

func ConvertSessionModelToRepoModel(session model.Session) repoModel.Session {
	return repoModel.Session{
		UserUUID:  session.UUID(),
		CreatedAt: session.CreatedAt().UnixNano(),
		UpdatedAt: session.UpdatedAt().UnixNano(),
		ExpiresAt: session.ExpiresAt().UnixNano(),
	}
}
