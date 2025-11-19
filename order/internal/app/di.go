package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	google_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/vipshark78/microservices-course-homeworks/order/internal/api/order/v1"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc"
	inventory "github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc/inventory/v1"
	payment "github.com/vipshark78/microservices-course-homeworks/order/internal/client/grpc/payment/v1"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository"
	orderRepository "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/order"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/service"
	orderService "github.com/vipshark78/microservices-course-homeworks/order/internal/service/order"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/migrator"
	orderMigrator "github.com/vipshark78/microservices-course-homeworks/platform/pkg/migrator/pg"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API          order_v1.Handler
	orderService        service.OrderService
	orderRepository     repository.OrderRepository
	pgxPool             *pgxpool.Pool
	inventoryClient     grpc.InventoryClient
	paymentClient       grpc.PaymentClient
	inventoryClientConn *google_grpc.ClientConn
	paymentClientConn   *google_grpc.ClientConn
	orderMigrator       migrator.Migrator
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) order_v1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx))
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PgxPool(ctx))

		d.OrderMigrator(d.PgxPool(ctx))
		err := d.orderMigrator.Up()
		if err != nil {
			logger.Fatal(ctx, "Ошибка при миграции", zap.Error(err))
		}
	}

	return d.orderRepository
}

func (d *diContainer) PgxPool(ctx context.Context) *pgxpool.Pool {
	if d.pgxPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			logger.Fatal(ctx, "failed to connect to database", zap.Error(err))
		}

		closer.AddNamed("PgxPool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgxPool = pool
	}

	return d.pgxPool
}

func (d *diContainer) OrderMigrator(pgxPool *pgxpool.Pool) migrator.Migrator {
	if d.orderMigrator == nil {
		d.orderMigrator = orderMigrator.NewMigrator(stdlib.OpenDB(*pgxPool.Config().ConnConfig), config.AppConfig().Postgres.MigrationsDir())
	}
	return d.orderMigrator
}

func (d *diContainer) InventoryClient(ctx context.Context) grpc.InventoryClient {
	if d.inventoryClient == nil {
		inventoryServiceClient := inventory_v1.NewInventoryServiceClient(d.InventoryClientConn(ctx))

		inventoryClient := inventory.NewInventoryClient(inventoryServiceClient)

		d.inventoryClient = inventoryClient
	}

	return d.inventoryClient
}

func (d *diContainer) InventoryClientConn(ctx context.Context) *google_grpc.ClientConn {
	if d.inventoryClientConn == nil {
		inventoryClientConn, err := google_grpc.NewClient(config.AppConfig().InventoryGRPC.Address(), google_grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal(ctx, "Ошибка подключения к сервису Inventory", zap.Error(err))
		}
		d.inventoryClientConn = inventoryClientConn

		closer.AddNamed("InventoryClientConnect", func(ctx context.Context) error {
			err = inventoryClientConn.Close()
			if err != nil {
				logger.Error(ctx, "Ошибка закрытия подключения к сервису Inventory", zap.Error(err))
			}
			return err
		})
	}

	return d.inventoryClientConn
}

func (d *diContainer) PaymentClient(ctx context.Context) grpc.PaymentClient {
	if d.paymentClient == nil {
		paymentServiceClient := payment_v1.NewPaymentServiceClient(d.PaymentClientConn(ctx))

		paymentClient := payment.NewPaymentClient(paymentServiceClient)

		d.paymentClient = paymentClient
	}

	return d.paymentClient
}

func (d *diContainer) PaymentClientConn(ctx context.Context) *google_grpc.ClientConn {
	if d.paymentClientConn == nil {
		paymentClientConn, err := google_grpc.NewClient(config.AppConfig().PaymentGRPC.Address(), google_grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal(ctx, "Ошибка подключения к сервису Payment", zap.Error(err))
		}
		d.paymentClientConn = paymentClientConn

		closer.AddNamed("PaymentClientConnect", func(ctx context.Context) error {
			err = paymentClientConn.Close()
			if err != nil {
				logger.Error(ctx, "Ошибка закрытия подключения к сервису Payment", zap.Error(err))
			}
			return err
		})
	}

	return d.paymentClientConn
}
