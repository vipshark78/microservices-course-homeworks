package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/payment/internal/app"
	"github.com/vipshark78/microservices-course-homeworks/payment/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func main() {
	ctx := context.Background()

	err := config.Load("../deploy/compose/payment/.env")
	if err != nil {
		panic(fmt.Errorf("ошибка при загрузке конфига Payment:%w", err))
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
	defer gracefulShutdown()

	paymentApp, err := app.New(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при создании приложения Payment", zap.Error(err))
	}

	err = paymentApp.Run(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при запуске приложения Payment", zap.Error(err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы Payment", zap.Error(err))
	}
}
