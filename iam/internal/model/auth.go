package model

import (
	"errors"
	"time"
)

type WhoamiResponse struct {
	Session Session
	User    User
}

type Session struct {
	uuid      string
	createdAt time.Time
	updatedAt time.Time
	expiresAt time.Time
}

func NewSession(userUUID string, ttl time.Duration) (Session, error) {
	if ttl < 0 {
		return Session{}, errors.New("ttl must be greater than zero")
	}

	return Session{
		uuid:      userUUID,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		expiresAt: time.Now().Add(ttl),
	}, nil
}

func NewSessionFromValues(uuid string, createdAt, updatedAt, expiresAt time.Time) Session {
	return Session{
		uuid:      uuid,
		createdAt: createdAt,
		updatedAt: updatedAt,
		expiresAt: expiresAt,
	}
}

func (s *Session) UUID() string {
	return s.uuid
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}
