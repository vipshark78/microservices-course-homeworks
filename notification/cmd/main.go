package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/notification/internal/app"
	"github.com/vipshark78/microservices-course-homeworks/notification/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

func main() {
	ctx := context.Background()

	err := config.Load("../deploy/compose/notification/.env")
	if err != nil {
		panic(fmt.Errorf("ошибка при загрузке конфига Notification:%w", err))
	}

	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	notificationApp, err := app.New(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при создании приложения Notification", zap.Error(err))
	}

	err = notificationApp.Run(ctx)
	if err != nil {
		logger.Fatal(ctx, "Ошибка при запуске приложения Notification", zap.Error(err))
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы Notification", zap.Error(err))
	}
}
