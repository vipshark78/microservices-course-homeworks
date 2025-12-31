package model

import "github.com/go-faster/errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalError      = errors.New("internal server error")
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidSessionUUID = errors.New("invalid session uuid")
	ErrInvalidUserUUID    = errors.New("invalid user uuid")
	ErrInvalidUserInfo    = errors.New("invalid user info")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrSessionExpired     = errors.New("session expired")
)
