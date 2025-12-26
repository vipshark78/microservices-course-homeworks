package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/converter"
	repomodel "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/model"
)

func (r *repository) GetByUUID(ctx context.Context, userUUID string) (model.User, error) {
	builderSelect := sq.
		Select(
			repomodel.ColumnUserUUID,
			repomodel.ColumnUserInfoLogin,
			repomodel.ColumnUserInfoEmail,
			repomodel.ColumnUserInfoPasswordHash,
			repomodel.ColumnCreatedAt,
			repomodel.ColumnUpdatedAt).
		From(repomodel.UsersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{repomodel.ColumnUserUUID: userUUID})

	return r.GetUser(ctx, builderSelect)
}

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	builderSelect := sq.
		Select(
			repomodel.ColumnUserUUID,
			repomodel.ColumnUserInfoLogin,
			repomodel.ColumnUserInfoEmail,
			repomodel.ColumnUserInfoPasswordHash,
			repomodel.ColumnCreatedAt,
			repomodel.ColumnUpdatedAt).
		From(repomodel.UsersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{repomodel.ColumnUserInfoLogin: login})

	user, err := r.GetUser(ctx, builderSelect)
	if err != nil {
		return model.User{}, err
	}

	notificationMethods, err := r.GetNotificationMethods(ctx, user.UserUUID)
	if err != nil {
		return model.User{}, err
	}

	user.UserInfo.NotificationMethods = converter.ConvertRepoNotificationMethodsToModel(notificationMethods)

	return user, nil
}

func (r *repository) GetNotificationMethods(ctx context.Context, userUUID string) ([]repomodel.NotificationMethod, error) {
	builderSelect := sq.
		Select(
			repomodel.ColumnNotificationProviderName,
			repomodel.ColumnNotificationTarget,
		).
		From(repomodel.NotificationTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{repomodel.ColumnUserUUID: userUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return []repomodel.NotificationMethod{}, err
	}

	var notificationMethods []repomodel.NotificationMethod

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to select notificationMethods: %v\n", err)
		return []repomodel.NotificationMethod{}, err
	}

	notificationMethods, err = pgx.CollectRows[repomodel.NotificationMethod](rows, pgx.RowToStructByName[repomodel.NotificationMethod])
	if err != nil {
		return []repomodel.NotificationMethod{}, fmt.Errorf("failed from collect notification methods: %w", err)
	}

	return notificationMethods, nil
}

func (r *repository) GetUser(ctx context.Context, builderSelect sq.SelectBuilder) (model.User, error) {
	var user repomodel.User
	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.User{}, err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&user.UserUUID,
		&user.UserInfo.Login,
		&user.UserInfo.Email,
		&user.UserInfo.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("failed to select user: %v\n", err)
		return model.User{}, err
	}

	return converter.ConvertRepoUserModelToModel(user), nil
}
