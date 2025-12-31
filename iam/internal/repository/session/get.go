package session

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/converter"
	repoModel "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUUID string) (model.Session, error) {
	values, err := r.cache.HGetAll(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.Session{}, model.ErrSessionNotFound
		}
		return model.Session{}, err
	}
	if len(values) == 0 {
		return model.Session{}, model.ErrSessionNotFound
	}

	var session repoModel.Session
	err = redigo.ScanStruct(values, &session)
	if err != nil {
		return model.Session{}, err
	}

	return converter.ConvertRepoSessionModelToModel(session), nil
}
