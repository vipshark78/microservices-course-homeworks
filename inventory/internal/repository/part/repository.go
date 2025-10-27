package part

import (
	"sync"

	repoModel "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
)

type repository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewRepository() *repository {
	repo := &repository{
		mu:    sync.RWMutex{},
		parts: make(map[string]repoModel.Part),
	}

	repo.initParts()

	return repo
}
