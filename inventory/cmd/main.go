package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/interceptor"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер с интерцептором логирования
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
		),
	)

	// Создаем хранилище деталей
	storage := newInventoryStorage()

	// Регистрируем наш сервис
	service := newInventoryService(storage)

	// Заполняем хранилище случайными данными
	service.initParts()

	inventory_v1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

// InventoryService сервис взаимодействия с деталями
type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer
	inventoryStorage *inventoryStorage
}

// newInventoryService создает новый экземпляр сервиса деталей
func newInventoryService(storage *inventoryStorage) *inventoryService {
	return &inventoryService{
		inventoryStorage: storage,
	}
}

// GetPart возвращает информацию о детали по UUID
func (i *inventoryService) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	part := i.inventoryStorage.GetPart(req.Uuid)
	if part == nil {
		return nil, status.Errorf(codes.NotFound, "part not found")
	}
	return &inventory_v1.GetPartResponse{Part: part}, nil
}

// ListParts возвращает список деталей с возможностью фильтрации
func (i *inventoryService) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	filter := req.Filter
	parts := make([]*inventory_v1.Part, 0, len(i.inventoryStorage.parts))
	filteredParts := i.filterParts(filter)
	parts = append(parts, filteredParts...)
	return &inventory_v1.ListPartsResponse{Parts: parts}, nil
}

// filterParts фильтрует детали
func (i *inventoryService) filterParts(filter *inventory_v1.PartsFilter) []*inventory_v1.Part {
	allParts := i.inventoryStorage.GetAllParts()
	if len(allParts) == 0 {
		return nil
	}
	if filter == nil {
		return allParts
	}
	if len(filter.Names) == 0 && len(filter.Categories) == 0 && len(filter.ManufacturerCountries) == 0 && len(filter.Uuids) == 0 && len(filter.Tags) == 0 {
		return allParts
	}
	filteredParts := make([]*inventory_v1.Part, 0, len(allParts))

	for _, part := range allParts {
		if len(filter.Uuids) > 0 && !slices.Contains(filter.Uuids, part.Uuid) {
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
		if len(filter.Tags) > 0 && !hasAnyTag(part.Tags, filter.Tags) {
			continue
		}
		filteredParts = append(filteredParts, part)
	}

	return filteredParts
}

// hasAnyTag проверяет наличие хотя бы одного тега из списка в информации о детали
func hasAnyTag(partTags, filterTags []string) bool {
	for _, tag := range filterTags {
		if slices.Contains(partTags, tag) {
			return true
		}
	}

	return false
}

// inventoryStorage хранилище деталей
type inventoryStorage struct {
	sync.RWMutex
	parts map[string]*inventory_v1.Part
}

// newInventoryStorage создает новое хранилище деталей
func newInventoryStorage() *inventoryStorage {
	return &inventoryStorage{parts: make(map[string]*inventory_v1.Part)}
}

// GetPart возвращает деталь по UUID
func (i *inventoryStorage) GetPart(uuid string) *inventory_v1.Part {
	i.RLock()
	defer i.RUnlock()
	part, ok := i.parts[uuid]
	if ok {
		return part
	}
	return nil
}

// GetAllParts возвращает все детали
func (i *inventoryStorage) GetAllParts() []*inventory_v1.Part {
	i.RLock()
	defer i.RUnlock()
	allParts := make([]*inventory_v1.Part, 0, len(i.parts))
	for _, part := range i.parts {
		allParts = append(allParts, part)
	}
	return allParts
}

func (i *inventoryService) initParts() {
	parts := generateParts()

	for _, part := range parts {
		i.inventoryStorage.parts[part.Uuid] = part
	}
}

func generateParts() []*inventory_v1.Part {
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

	var parts []*inventory_v1.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &inventory_v1.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventory_v1.Category(gofakeit.Number(1, 4)), //nolint:gosec // safe: gofakeit.Number returns 1..4
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		})
	}

	return parts
}

func generateDimensions() *inventory_v1.Dimensions {
	return &inventory_v1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventory_v1.Manufacturer {
	return &inventory_v1.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]*inventory_v1.Value {
	metadata := make(map[string]*inventory_v1.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

func generateMetadataValue() *inventory_v1.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventory_v1.Value{
			ValueType: &inventory_v1.Value_StringValue{
				StringValue: gofakeit.Word(),
			},
		}

	case 1:
		return &inventory_v1.Value{
			ValueType: &inventory_v1.Value_Int64Value{
				Int64Value: int64(gofakeit.Number(1, 100)),
			},
		}

	case 2:
		return &inventory_v1.Value{
			ValueType: &inventory_v1.Value_DoubleValue{
				DoubleValue: roundTo(gofakeit.Float64Range(1, 100)),
			},
		}

	case 3:
		return &inventory_v1.Value{
			ValueType: &inventory_v1.Value_BooleanValue{
				BooleanValue: gofakeit.Bool(),
			},
		}

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
