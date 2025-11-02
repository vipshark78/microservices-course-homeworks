package part

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *repository {
	repo := &repository{
		db: db,
	}

	return repo
}
