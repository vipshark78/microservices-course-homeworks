package converter

import (
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

// ProtoPartToModelParts конвертирует protobuf Parts в модель Parts.
func ProtoPartsToModelParts(parts []*inventory_v1.Part) []model.Part {
	result := make([]model.Part, len(parts))
	for _, part := range parts {
		result = append(result, protoPartToModelPart(part))
	}
	return result
}

// ModelToPartsFilter конвертирует модель PartsFilter в protobuf PartsFilter.
func ModelToPartsFilter(filter model.PartsFilter) *inventory_v1.PartsFilter {
	return &inventory_v1.PartsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            modelCategoriesToCategories(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// modelCategoriesToCategories конвертирует модель Categories в protobuf Categories.
func modelCategoriesToCategories(categories []string) []inventory_v1.Category {
	result := make([]inventory_v1.Category, len(categories))
	for i, category := range categories {
		switch category {
		case model.PORTHOLE:

			result[i] = inventory_v1.Category_PORTHOLE
		case model.ENGINE:
			result[i] = inventory_v1.Category_ENGINE
		case model.FUEL:
			result[i] = inventory_v1.Category_FUEL
		case model.WING:
			result[i] = inventory_v1.Category_WING
		default:
			result[i] = inventory_v1.Category_UNKNOWN_UNSPECIFIED
		}
	}
	return result
}

// protoPartToModelPart конвертирует protobuf Part в модель Part.
func protoPartToModelPart(part *inventory_v1.Part) model.Part {
	return model.Part{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      protoCategoryToModelCategory(part.Category),
		Dimensions:    protoDimensionsToModelDimensions(part.Dimensions),
		Manufacturer:  protoManufacturerToModelManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      protoMetadataToModelMetadata(part.Metadata),
	}
}

// protoCategoryToModelCategory конвертирует protobuf Category в модель Category.
func protoCategoryToModelCategory(category inventory_v1.Category) string {
	switch category {
	case inventory_v1.Category_PORTHOLE:
		return model.PORTHOLE
	case inventory_v1.Category_ENGINE:
		return model.ENGINE
	case inventory_v1.Category_FUEL:
		return model.FUEL
	case inventory_v1.Category_WING:
		return model.WING
	default:
		return model.UNKNOWN_UNSPECIFIED
	}
}

// protoDimensionsToModelDimensions конвертирует protobuf Dimensions в модель Dimensions.
func protoDimensionsToModelDimensions(dimensions *inventory_v1.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// protoManufacturerToModelManufacturer конвертирует protobuf Manufacturer в модель Manufacturer.
func protoManufacturerToModelManufacturer(manufacturer *inventory_v1.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// protoMetadataToModelMetadata конвертирует protobuf Metadata в модель Metadata.
func protoMetadataToModelMetadata(metadata map[string]*inventory_v1.Value) map[string]*model.Value {
	result := make(map[string]*model.Value, len(metadata))
	for k, v := range metadata {
		result[k] = protoValueToModelValue(v)
	}
	return result
}

// protoValueToModelValue конвертирует protobuf Value в модель Value.
func protoValueToModelValue(value *inventory_v1.Value) *model.Value {
	if value == nil {
		return nil
	}
	switch v := value.ValueType.(type) {
	case *inventory_v1.Value_StringValue:
		return &model.Value{StringValue: &v.StringValue}
	case *inventory_v1.Value_Int64Value:
		return &model.Value{Int64Value: &v.Int64Value}
	case *inventory_v1.Value_DoubleValue:
		return &model.Value{DoubleValue: &v.DoubleValue}
	case *inventory_v1.Value_BooleanValue:
		return &model.Value{BooleanValue: &v.BooleanValue}
	default:
		return nil
	}
}
