package part

import (
	"math"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
)

// initParts инициализирует репозиторий случайными данными.
func (r *repository) initParts() {
	parts := generateParts()

	for _, part := range parts {
		r.parts[part.UUID] = part
	}
}

// generateParts генерирует список случайных данных для заполнения репозитория.
func generateParts() []model.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	var parts []model.Part

	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		createdTime := time.Now()
		parts = append(parts, model.Part{
			UUID:          gofakeit.UUID(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      gofakeit.RandomString([]string{model.FUEL, model.WING, model.ENGINE, model.PORTHOLE}),
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     &createdTime,
			UpdatedAt:     &createdTime,
		})
	}

	return parts
}

// generateDimensions генерирует случайные размеры.
func generateDimensions() *model.Dimensions {
	return &model.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

// generateManufacturer генерирует случайную информацию о производителе.
func generateManufacturer() *model.Manufacturer {
	return &model.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

// generateTags генерирует случайный набор тегов.
func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

// generateMetadata генерирует случайное количество метаданных.
func generateMetadata() map[string]*model.Value {
	metadata := make(map[string]*model.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

// generateMetadataValue генерирует случайное значение метаданных.
func generateMetadataValue() *model.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &model.Value{
			StringValue: lo.ToPtr(gofakeit.Word()),
		}

	case 1:
		return &model.Value{
			Int64Value: lo.ToPtr(int64(gofakeit.Number(1, 100))),
		}

	case 2:
		return &model.Value{
			DoubleValue: lo.ToPtr(gofakeit.Float64Range(1, 100)),
		}

	case 3:
		return &model.Value{
			BooleanValue: lo.ToPtr(gofakeit.Bool()),
		}

	default:
		return nil
	}
}

// roundTo округляет число до двух знаков после запятой.
func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
