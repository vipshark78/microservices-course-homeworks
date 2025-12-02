package model

import "time"

type User struct {
	UserUUID  string
	UserInfo  UserInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UserInfo struct {
	Login               string
	Email               string
	PasswordHash        string
	NotificationMethods []NotificationMethod
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type UserRegistrationInfo struct {
	Info     UserInfo
	Password string
}
