package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpcMiddleware "github.com/vipshark78/microservices-course-homeworks/platform/pkg/middleware/grpc"
	repoModel "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/model"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет тестовую деталь в коллекцию Mongo и возвращает ее UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	part := repoModel.Part{
		UUID:          partUUID,
		Name:          gofakeit.Name(),
		Description:   gofakeit.Word(),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 1000)),
		Category:      repoModel.ENGINE,
		Dimensions: &repoModel.Dimensions{
			Length: gofakeit.Float64Range(1, 1000),
			Weight: gofakeit.Float64Range(1, 1000),
			Height: gofakeit.Float64Range(1, 1000),
			Width:  gofakeit.Float64Range(1, 1000),
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{gofakeit.EmojiTag(), gofakeit.EmojiTag(), gofakeit.EmojiTag(), gofakeit.EmojiTag()},
		Metadata:  map[string]repoModel.Value{gofakeit.Word(): {StringValue: lo.ToPtr(gofakeit.Word())}},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовую деталь с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part repoModel.Part) (string, error) {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return part.UUID, nil
}

// GetTestPartInfo — возвращает тестовую информацию о детали
func (env *TestEnvironment) GetTestPartInfo() *inventory_v1.Part {
	return &inventory_v1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Word(),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 1000)),
		Category:      inventory_v1.Category_ENGINE,
		Dimensions: &inventory_v1.Dimensions{
			Length: gofakeit.Float64Range(1, 1000),
			Width:  gofakeit.Float64Range(1, 1000),
			Weight: gofakeit.Float64Range(1, 1000),
			Height: gofakeit.Float64Range(1, 1000),
		},
		Manufacturer: &inventory_v1.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{gofakeit.EmojiTag(), gofakeit.EmojiTag(), gofakeit.EmojiTag()},
		Metadata:  map[string]*inventory_v1.Value{gofakeit.Word(): {ValueType: &inventory_v1.Value_StringValue{StringValue: gofakeit.Word()}}},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

// AddSessionUUIDToContext добавляет session-uuid в исходящие gRPC metadata контекста
func AddSessionUUIDToContext(ctx context.Context, sessionUUID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, grpcMiddleware.SessionUUIDMetadataKey, sessionUUID)
}
