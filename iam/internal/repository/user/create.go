package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/converter"
	repomodel "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, info model.UserInfo, hash string) (string, error) {
	repoInfo := converter.ConvertUserInfoToRepoModel(info, hash)
	builderInsert := sq.
		Insert(repomodel.TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			repomodel.ColumnUserInfoLogin,
			repomodel.ColumnUserInfoEmail,
			repomodel.ColumnUserInfoPasswordHash,
			repomodel.ColumnUserInfoNotificationMethods,
		).
		Values(
			repoInfo.Login,
			repoInfo.Email,
			repoInfo.PasswordHash,
			repoInfo.NotificationMethods,
		).
		Suffix(fmt.Sprintf("RETURNING %s", repomodel.ColumnUserUUID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return "", err
	}

	var userUUID string
	err = r.pool.QueryRow(ctx, query, args...).Scan(&userUUID)
	if err != nil {
		log.Printf("failed to insert user: %v\n", err)
		return "", err
	}
	return userUUID, nil
}
