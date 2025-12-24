package app

import (
	"context"
	"fmt"

	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/metrics"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
	platformMetrics "github.com/vipshark78/microservices-course-homeworks/platform/pkg/metrics"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runConsumer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMetrics,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	err := logger.Init(
		ctx,
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
		config.AppConfig().Logger.EnableOTLP(),
		config.AppConfig().Logger.OtelCollectorEndpoint(),
		config.AppConfig().Logger.ServiceName(),
	)
	if err != nil {
		return err
	}

	closer.AddNamed("Logger", func(ctx context.Context) error {
		err = logger.CloseAndSync(ctx)
		if err != nil {
			return fmt.Errorf("ошибка при закрытии логгера:%w", err)
		}
		return nil
	})
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	err := platformMetrics.InitProvider(ctx, config.AppConfig().Metrics)
	if err != nil {
		return err
	}
	err = metrics.InitMetrics()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer запущен")

	err := a.diContainer.ConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}
	return nil
}
