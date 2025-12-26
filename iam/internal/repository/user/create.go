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
		Insert(repomodel.UsersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			repomodel.ColumnUserInfoLogin,
			repomodel.ColumnUserInfoEmail,
			repomodel.ColumnUserInfoPasswordHash,
		).
		Values(
			repoInfo.Login,
			repoInfo.Email,
			repoInfo.PasswordHash,
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

	if len(repoInfo.NotificationMethods) > 0 {
		builderInsert = sq.
			Insert(repomodel.NotificationTableName).
			PlaceholderFormat(sq.Dollar).
			Columns(
				repomodel.ColumnUserUUID,
				repomodel.ColumnNotificationProviderName,
				repomodel.ColumnNotificationTarget,
			)

		for _, m := range repoInfo.NotificationMethods {
			builderInsert = builderInsert.Values(userUUID, m.ProviderName, m.Target)
		}

		query, args, err = builderInsert.ToSql()
		if err != nil {
			log.Printf("failed to build query: %v\n", err)
			return "", err
		}

		_, err = r.pool.Exec(ctx, query, args...)
		if err != nil {
			log.Printf("failed to insert notification methods: %v\n", err)
			return "", err
		}
	}

	return userUUID, nil
}
