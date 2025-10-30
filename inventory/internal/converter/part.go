package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/model"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

// ModelsToParts конвертирует модели бизнес-логики в прото-модели.
func ModelsToParts(parts []model.Part) []*inventory_v1.Part {
	result := make([]*inventory_v1.Part, len(parts))
	for i, part := range parts {
		result[i] = ModelToPart(part)
	}
	return result
}

// ModelToPart конвертирует модель бизнес-логики в прото-модель.
func ModelToPart(part model.Part) *inventory_v1.Part {
	return &inventory_v1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventory_v1.Category(inventory_v1.Category_value[part.Category]),
		Dimensions:    DimensionsToProto(part.Dimensions),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(*part.CreatedAt),
		UpdatedAt:     timestamppb.New(*part.UpdatedAt),
	}
}

// DimensionsToProto конвертирует модель бизнес-логики размеров в прото-модель.
func DimensionsToProto(dimensions *model.Dimensions) *inventory_v1.Dimensions {
	return &inventory_v1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// ManufacturerToProto конвертирует модель бизнес-логики производителя в прото-модель.
func ManufacturerToProto(manufacturer *model.Manufacturer) *inventory_v1.Manufacturer {
	return &inventory_v1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// MetadataToProto конвертирует модель бизнес-логики метаданных в прото-модель.
func MetadataToProto(metadata map[string]*model.Value) map[string]*inventory_v1.Value {
	result := make(map[string]*inventory_v1.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ValueToProto(value)
	}
	return result
}

// ValueToProto конвертирует модель бизнес-логики значения в прото-модель.
func ValueToProto(value *model.Value) *inventory_v1.Value {
	if value == nil {
		return nil
	}
	if value.StringValue != nil {
		return &inventory_v1.Value{ValueType: &inventory_v1.Value_StringValue{StringValue: *value.StringValue}}
	}
	if value.Int64Value != nil {
		return &inventory_v1.Value{ValueType: &inventory_v1.Value_Int64Value{Int64Value: *value.Int64Value}}
	}
	if value.DoubleValue != nil {
		return &inventory_v1.Value{ValueType: &inventory_v1.Value_DoubleValue{DoubleValue: *value.DoubleValue}}
	}
	if value.BooleanValue != nil {
		return &inventory_v1.Value{ValueType: &inventory_v1.Value_BooleanValue{BooleanValue: *value.BooleanValue}}
	}
	return nil
}

// DimensionsToModel конвертирует прото-модель размеров в модель бизнес-логики.
func DimensionsToModel(dimensions *inventory_v1.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// ManufacturerToModel конвертирует прото-модель производителя в модель бизнес-логики.
func ManufacturerToModel(manufacturer *inventory_v1.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// MetadataToModel конвертирует прото-модель метаданных в модель бизнес-логики.
func MetadataToModel(metadata map[string]*inventory_v1.Value) map[string]*model.Value {
	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ValueToModel(value)
	}
	return result
}

// ValueToModel конвертирует прото-модель значения в модель бизнес-логики.
func ValueToModel(value *inventory_v1.Value) *model.Value {
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

// PartsFilterToModel конвертирует прото-модель фильтра в модель бизнес-логики.
func PartsFilterToModel(filter *inventory_v1.PartsFilter) model.PartsFilter {
	return model.PartsFilter{
		UUIDs:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoriesToModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// CategoriesToModel конвертирует прото-модель категорий в модель бизнес-логики.
func CategoriesToModel(categories []inventory_v1.Category) []string {
	result := make([]string, len(categories))
	for i, category := range categories {
		result[i] = category.String()
	}
	return result
}
