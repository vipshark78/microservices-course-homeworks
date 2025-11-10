package part

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/converter"
	repoModel "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
)

// GetPart получение детали по UUID.
func (r *repository) GetPart(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	c := r.db.Collection(repoModel.CollectionName)
	res := c.FindOne(ctx, bson.M{repoModel.CollectionFieldUUID: uuid.String()})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}
	var repoPart repoModel.Part
	err := res.Decode(&repoPart)
	if err != nil {
		return model.Part{}, err
	}
	return converter.ModelToPart(repoPart), nil
}
