package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

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
			repomodel.ColumnUserInfoNotificationMethods,
			repomodel.ColumnCreatedAt,
			repomodel.ColumnUpdatedAt).
		From(repomodel.TableName).
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
			repomodel.ColumnUserInfoNotificationMethods,
			repomodel.ColumnCreatedAt,
			repomodel.ColumnUpdatedAt).
		From(repomodel.TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{repomodel.ColumnUserInfoLogin: login})

	return r.GetUser(ctx, builderSelect)
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
		&user.UserInfo.NotificationMethods,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("failed to select user: %v\n", err)
		return model.User{}, err
	}

	return converter.ConvertRepoUserModelToModel(user), nil
}
