package part

import (
	"context"

	repoModel "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
)

func (r *repository) InsertPart(ctx context.Context, p repoModel.Part) error {
	c := r.db.Collection(repoModel.CollectionName)
	_, err := c.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}
