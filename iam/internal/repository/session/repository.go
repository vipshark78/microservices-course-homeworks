package session

import (
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/cache"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}
