package model

type Session struct {
	UserUUID  string `redis:"user_uuid"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt int64  `redis:"updated_at,omitempty"`
	ExpiresAt int64  `redis:"expires_at"`
}
