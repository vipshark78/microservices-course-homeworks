package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/app"
	"github.com/vipshark78/microservices-course-homeworks/inventory/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func main() {
	ctx := context.Background()

	err := config.Load("../deploy/compose/inventory/.env")
	if err != nil {
		panic(fmt.Errorf("ошибка при загрузке конфига Inventory:%w", err))
	}

	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	inventoryApp, err := app.New(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при создании приложения Inventory", zap.Error(err))
	}

	err = inventoryApp.Run(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при запуске приложения Inventory", zap.Error(err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы Inventory", zap.Error(err))
	}
}
