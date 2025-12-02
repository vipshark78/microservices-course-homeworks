package converter

import (
	"time"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	repoModel "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/model"
)

func ConvertRepoSessionModelToModel(repoSession repoModel.Session) model.Session {
	return model.Session{
		UUID:      repoSession.UserUUID,
		CreatedAt: time.Unix(repoSession.CreatedAt, 0),
		UpdatedAt: time.Unix(repoSession.UpdatedAt, 0),
		ExpiresAt: time.Unix(repoSession.ExpiresAt, 0),
	}
}
