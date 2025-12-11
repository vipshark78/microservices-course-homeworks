package converter

import (
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
	common_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/common/v1"
	user_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/user/v1"
)

func ConvertServiceUserToProtoModel(user model.User) *common_v1.User {
	return &common_v1.User{
		Uuid: user.UserUUID,
		Info: ConvertServiceUserInfoToProtoModel(user.UserInfo),
	}
}

func ConvertServiceUserInfoToProtoModel(userInfo model.UserInfo) *common_v1.UserInfo {
	return &common_v1.UserInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		NotificationMethods: ConvertServiceNotificationMethodsToProtoModel(userInfo.NotificationMethods),
	}
}

func ConvertServiceNotificationMethodsToProtoModel(notificationMethods []model.NotificationMethod) []*common_v1.NotificationMethod {
	methods := make([]*common_v1.NotificationMethod, len(notificationMethods))
	for _, method := range notificationMethods {
		methods = append(methods, &common_v1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}
	return methods
}

func ConvertProtoUserRegInfoToServiceModel(regInfo *user_v1.UserRegistrationInfo) model.UserRegistrationInfo {
	return model.UserRegistrationInfo{
		Info:     ConvertProtoUserInfoToServiceModel(regInfo.Info),
		Password: regInfo.GetPassword(),
	}
}

func ConvertProtoUserInfoToServiceModel(info *common_v1.UserInfo) model.UserInfo {
	return model.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		PasswordHash:        "",
		NotificationMethods: ConvertProtoNotificationMethodsToServiceModel(info.NotificationMethods),
	}
}

func ConvertProtoNotificationMethodsToServiceModel(methods []*common_v1.NotificationMethod) []model.NotificationMethod {
	var result []model.NotificationMethod
	for _, method := range methods {
		result = append(result, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}
	return result
}
