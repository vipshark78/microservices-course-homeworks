package converter

import (
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	repomodel "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/model"
)

func ConvertUserInfoToRepoModel(info model.UserInfo, hash string) repomodel.UserInfo {
	return repomodel.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		PasswordHash:        hash,
		NotificationMethods: ConvertNotificationMethodsToRepoModel(info.NotificationMethods),
	}
}

func ConvertNotificationMethodsToRepoModel(notificationMethods []model.NotificationMethod) []repomodel.NotificationMethod {
	var methods []repomodel.NotificationMethod
	for _, method := range notificationMethods {
		methods = append(methods, repomodel.NotificationMethod{ProviderName: method.ProviderName, Target: method.Target})
	}
	return methods
}

func ConvertRepoUserModelToModel(repoUser repomodel.User) model.User {
	return model.User{
		UserUUID:  repoUser.UserUUID,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
		UserInfo: model.UserInfo{
			Login:               repoUser.UserInfo.Login,
			Email:               repoUser.UserInfo.Email,
			PasswordHash:        repoUser.UserInfo.PasswordHash,
			NotificationMethods: ConvertRepoNotificationMethodsToModel(repoUser.UserInfo.NotificationMethods),
		},
	}
}

func ConvertRepoNotificationMethodsToModel(notificationMethods []repomodel.NotificationMethod) []model.NotificationMethod {
	var methods []model.NotificationMethod
	for _, method := range notificationMethods {
		methods = append(methods, model.NotificationMethod{ProviderName: method.ProviderName, Target: method.Target})
	}
	return methods
}
