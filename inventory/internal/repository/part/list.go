package part

import (
	"context"
	"slices"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/converter"
)

// ListParts возвращает список деталей по заданному фильтру. Если фильтр пустой - вернет все детали.
func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.filterParts(ctx, filter)
}

// filterParts фильтрует детали
func (r *repository) filterParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	allParts := make([]model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		allParts = append(allParts, converter.ModelToPart(part))
	}

	if r.isAllFiltersEmpty(filter) {
		return allParts, nil
	}
	filteredParts := make([]model.Part, 0, len(allParts))

	for _, part := range allParts {
		if len(filter.UUIDs) > 0 && !slices.Contains(filter.UUIDs, part.UUID) {
			continue
		}
		if len(filter.Names) > 0 && !slices.Contains(filter.Names, part.Name) {
			continue
		}
		if len(filter.Categories) > 0 && !slices.Contains(filter.Categories, part.Category) {
			continue
		}
		if len(filter.ManufacturerCountries) > 0 {
			if part.Manufacturer == nil {
				continue
			}
			if !slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
				continue
			}
		}
		if len(filter.Tags) > 0 && !r.hasAnyTag(part.Tags, filter.Tags) {
			continue
		}
		filteredParts = append(filteredParts, part)
	}

	if len(filteredParts) == 0 {
		return nil, model.ErrPartNotFound
	}

	return filteredParts, nil
}

// isAllFiltersEmpty проверяет, что все фильтры пустые
func (r *repository) isAllFiltersEmpty(filter model.PartsFilter) bool {
	if len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.UUIDs) == 0 &&
		len(filter.Tags) == 0 {
		return true
	}
	return false
}

// hasAnyTag проверяет наличие хотя бы одного тега из списка в информации о детали
func (r *repository) hasAnyTag(partTags, filterTags []string) bool {
	for _, tag := range filterTags {
		if slices.Contains(partTags, tag) {
			return true
		}
	}

	return false
}
