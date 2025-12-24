package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/config"
	orderMetrics "github.com/vipshark78/microservices-course-homeworks/order/internal/metrics"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
	platformMetrics "github.com/vipshark78/microservices-course-homeworks/platform/pkg/metrics"
	httpMidlleware "github.com/vipshark78/microservices-course-homeworks/platform/pkg/middleware/http"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/tracing"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Запускаем консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer error: %w", err)
		}
	}()

	// Запускаем HTTP-сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		logger.Error(ctx, "❌ Компонент завершился с ошибкой, завершение работы", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	case <-ctx.Done():
		logger.Info(ctx, "🔔 Получен сигнал завершения работы")
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initHTTPServer,
		a.initTracing,
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

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
	if err != nil {
		return err
	}

	closer.AddNamed("tracer", tracing.ShutdownTracer)

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	// Инициализируем платформенный провайдер метрик
	// Это создает MeterProvider и настраивает отправку в OTLP Collector
	err := platformMetrics.InitProvider(ctx, config.AppConfig().Metrics)
	if err != nil {
		return fmt.Errorf("failed to init metrics provider: %w", err)
	}

	// Инициализируем метрики Order сервиса
	// Это создает конкретные метрики (OrdersTotal и OrdersRevenueTotal) через Meter
	err = orderMetrics.InitMetrics()
	if err != nil {
		return fmt.Errorf("failed to init order metrics: %w", err)
	}

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
			err = fmt.Errorf("ошибка при закрытии логгера:%w", err)
		}
		return err
	})
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	authMiddleware := httpMidlleware.NewAuthMiddleware(a.diContainer.IAMClient(ctx))
	orderServer, err := order_v1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return err
	}

	// Создаем маршрутизатор HTTP
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(config.AppConfig().OrderHTTP.OperationTimeout()))
	router.Use(authMiddleware.Handle)

	// Монтируем обработчики OpenAPI
	router.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           router,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	closer.AddNamed("OrderHTTP", func(ctx context.Context) error {
		err = a.httpServer.Close()
		if err != nil {
			logger.Error(ctx, "Ошибка закрытия сервера OrderHTTP", zap.Error(err))
		}

		return err
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP-сервер запущен на порту %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Order Kafka consumer запущен")

	err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
