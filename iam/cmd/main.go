package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/app"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func main() {
	ctx := context.Background()

	err := config.Load("../deploy/compose/iam/.env")
	if err != nil {
		panic(fmt.Errorf("ошибка при загрузке конфига IAM:%w", err))
	}

	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	assemblyApp, err := app.New(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при создании приложения IAM", zap.Error(err))
	}

	err = assemblyApp.Run(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при запуске приложения IAM", zap.Error(err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы IAM", zap.Error(err))
	}
}
