package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/app"
	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func main() {
	ctx := context.Background()

	err := config.Load("../deploy/compose/assembly/.env")
	if err != nil {
		panic(fmt.Errorf("ошибка при загрузке конфига Assembly:%w", err))
	}

	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	assemblyApp, err := app.New(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при создании приложения Assembly", zap.Error(err))
	}

	err = assemblyApp.Run(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при запуске приложения Assembly", zap.Error(err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы Assembly", zap.Error(err))
	}
}
