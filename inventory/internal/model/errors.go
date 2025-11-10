package model

import (
	"errors"

	sharedErrors "github.com/vipshark78/microservices-course-homeworks/shared/pkg/errors"
)

var (
	ErrPartNotFound   = sharedErrors.NewNotFoundError(errors.New("part not found"))
	ErrInvalidUUID    = sharedErrors.NewInvalidArgumentError(errors.New("invalid uuid"))
	ErrInvalidRequest = sharedErrors.NewInvalidArgumentError(errors.New("invalid request"))
)
