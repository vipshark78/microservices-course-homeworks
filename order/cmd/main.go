package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/app"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func main() {
	ctx := context.Background()

	err := config.Load("../deploy/compose/order/.env")
	if err != nil {
		panic(fmt.Errorf("ошибка при загрузке конфига Order:%w", err))
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	defer gracefulShutdown()

	orderApp, err := app.New(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при создании приложения Order", zap.Error(err))
	}

	err = orderApp.Run(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при запуске приложения Order", zap.Error(err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы Order", zap.Error(err))
	}
}
