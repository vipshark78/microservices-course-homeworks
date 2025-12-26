package session

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/converter"
)

const ttl = time.Hour * 24

func (r *repository) Create(ctx context.Context, session model.Session) (string, error) {
	repoSession := converter.ConvertSessionModelToRepoModel(session)

	sessionUUID := uuid.New().String()

	err := r.cache.HashSet(ctx, sessionUUID, repoSession)
	if err != nil {
		return "", err
	}

	err = r.cache.Expire(ctx, sessionUUID, ttl)
	if err != nil {
		return "", err
	}

	err = r.AddToUserSet(ctx, repoSession.UserUUID, sessionUUID)
	if err != nil {
		return "", err
	}

	return sessionUUID, nil
}
