package session

import (
	"context"
	"time"

	"github.com/google/uuid"

	repoModel "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/model"
)

const ttl = time.Hour * 24

func (r *repository) Create(ctx context.Context, userUUID string) (string, error) {
	session := repoModel.Session{
		UserUUID:  userUUID,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: time.Now().Add(ttl).Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	sessionUUID := uuid.New().String()

	err := r.cache.HashSet(ctx, sessionUUID, session)
	if err != nil {
		return "", err
	}

	err = r.cache.Expire(ctx, sessionUUID, ttl)
	if err != nil {
		return "", err
	}

	err = r.AddToUserSet(ctx, userUUID, sessionUUID)
	if err != nil {
		return "", err
	}

	return sessionUUID, nil
}
