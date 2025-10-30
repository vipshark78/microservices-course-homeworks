package part

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/converter"
)

// GetPart получение детали по UUID.
func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	repoPart, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}
	return converter.ModelToPart(repoPart), nil
}
