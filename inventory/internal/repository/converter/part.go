package converter

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	repoModel "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
)

func FilterToBsonFilter(repoFilter model.PartsFilter) bson.M {
	filter := bson.M{}
	if len(repoFilter.UUIDs) > 0 {
		filter["uuid"] = bson.M{"$in": repoFilter.UUIDs}
	}
	if len(repoFilter.Names) > 0 {
		filter["name"] = bson.M{"$in": repoFilter.Names}
	}
	if len(repoFilter.Categories) > 0 {
		categories := make([]interface{}, len(repoFilter.Categories))
		for i, cat := range repoFilter.Categories {
			categories[i] = string(cat)
		}
		filter["category"] = bson.M{"$in": categories}
	}
	if len(repoFilter.ManufacturerCountries) > 0 {
		filter["manufacturer.country"] = bson.M{"$in": repoFilter.ManufacturerCountries}
	}
	if len(repoFilter.Tags) > 0 {
		filter["tags"] = bson.M{"$all": repoFilter.Tags}
	}
	return filter
}

// ModelToPart конвертирует модель репозитория в модель бизнес-логики.
func ModelToPart(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      ModelToCategory(part.Category),
		Dimensions:    ModelToDimensions(part.Dimensions),
		Manufacturer:  ModelToManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ModelToMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// ModelToCategory конвертирует модель репозитория в модель бизнес-логики.
func ModelToCategory(category repoModel.Category) model.Category {
	switch category {
	case repoModel.ENGINE:
		return model.ENGINE
	case repoModel.FUEL:
		return model.FUEL
	case repoModel.PORTHOLE:
		return model.PORTHOLE
	case repoModel.WING:
		return model.WING
	}
	return ""
}

// ModelToDimensions конвертирует модель репозитория в модель бизнес-логики.
func ModelToDimensions(dimensions *repoModel.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// ModelToManufacturer конвертирует модель репозитория в модель бизнес-логики.
func ModelToManufacturer(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// ModelToMetadata конвертирует модель репозитория в модель бизнес-логики.
func ModelToMetadata(metadata map[string]repoModel.Value) map[string]model.Value {
	result := make(map[string]model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ModelToValue(value)
	}
	return result
}

// ValueToModel конвертирует модель бизнес-логики в модель репозитория.
func ValueToModel(value *model.Value) *repoModel.Value {
	if value == nil {
		return nil
	}
	if value.StringValue != nil {
		return &repoModel.Value{StringValue: value.StringValue}
	}
	if value.Int64Value != nil {
		return &repoModel.Value{Int64Value: value.Int64Value}
	}
	if value.DoubleValue != nil {
		return &repoModel.Value{DoubleValue: value.DoubleValue}
	}
	if value.BooleanValue != nil {
		return &repoModel.Value{BooleanValue: value.BooleanValue}
	}
	return nil
}

// ModelToValue конвертирует модель репозитория в модель бизнес-логики.
func ModelToValue(value repoModel.Value) model.Value {
	if value.StringValue != nil {
		return model.Value{StringValue: value.StringValue}
	}
	if value.Int64Value != nil {
		return model.Value{Int64Value: value.Int64Value}
	}
	if value.DoubleValue != nil {
		return model.Value{DoubleValue: value.DoubleValue}
	}
	if value.BooleanValue != nil {
		return model.Value{BooleanValue: value.BooleanValue}
	}
	return model.Value{}
}
