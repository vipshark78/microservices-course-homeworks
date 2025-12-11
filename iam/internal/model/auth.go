package model

import "time"

type WhoamiResponse struct {
	Session Session
	User    User
}

type Session struct {
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
