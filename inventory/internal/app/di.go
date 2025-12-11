package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1API "github.com/vipshark78/microservices-course-homeworks/inventory/internal/api/inventory/v1"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository"
	inventoryRepository "github.com/vipshark78/microservices-course-homeworks/inventory/internal/repository/part"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/service"
	inventoryService "github.com/vipshark78/microservices-course-homeworks/inventory/internal/service/part"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
	grpcMidlleware "github.com/vipshark78/microservices-course-homeworks/platform/pkg/middleware/grpc"
	auth_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API      inventory_v1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository
	mongoDBClient       *mongo.Client
	mongoDBHandle       *mongo.Database
	iamClient           grpcMidlleware.IAMClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventory_v1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.PartRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(d.MongoDBHandle(ctx))
		err := d.inventoryRepository.Init(ctx)
		if err != nil {
			panic(fmt.Sprintf("Ошибка при инициализации репозитория: %v\n", err))
		}
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) IAMClient(ctx context.Context) grpcMidlleware.IAMClient {
	if d.iamClient == nil {
		d.iamClient = auth_v1.NewAuthServiceClient(d.IAMConn(ctx))
	}
	return d.iamClient
}

func (d *diContainer) IAMConn(_ context.Context) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		config.AppConfig().IAMGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("❌ Ошибка подключения к IAM Service: %v", err))
	}

	closer.AddNamed("IAM client", func(ctx context.Context) error {
		if err := conn.Close(); err != nil {
			logger.Error(ctx, "❌ Ошибка при закрытии подключения с IAM Service", zap.Error(err))
			return err
		}
		return nil
	})

	return conn
}
