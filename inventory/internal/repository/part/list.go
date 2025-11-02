package part

import (
	"context"
	"log"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/converter"
	repoModel "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
)

// ListParts возвращает список деталей по заданному фильтру. Если фильтр пустой - вернет все детали.
func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	return r.filterParts(ctx, filter)
}

// filterParts фильтрует детали
func (r *repository) filterParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	allParts := make([]model.Part, 0, 10)
	collection := r.db.Collection(repoModel.CollectionName)
	bsonFilter := converter.FilterToBsonFilter(filter)
	cursor, err := collection.Find(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = cursor.Close(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	for cursor.Next(ctx) {
		var part repoModel.Part
		err = cursor.Decode(&part)
		if err != nil {
			return nil, err
		}
		allParts = append(allParts, converter.ModelToPart(part))
	}
	return allParts, nil
}
