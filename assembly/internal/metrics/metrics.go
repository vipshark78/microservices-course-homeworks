package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	namespace   = "rocketFactory" // Пространство имен для группировки метрик
	serviceName = "assembly-service"
)

// =============================================================================
// METER - ФАБРИКА ДЛЯ СОЗДАНИЯ МЕТРИК
// =============================================================================
//
// Meter в OpenTelemetry - это фабрика для создания инструментов измерения метрик.
// Каждый КОМПОНЕНТ должен иметь свой meter с уникальным именем.
//
// АРХИТЕКТУРА ВЗАИМОДЕЙСТВИЯ:
//
//  1. platform/metrics инициализирует MeterProvider:
//     platform.InitProvider() → otel.SetMeterProvider(meterProvider)
//
//  2. assembly/metrics создает свой Meter:
//     otel.Meter("assembly-service") → получает глобальный MeterProvider
//
//  3. Meter создает метрики через MeterProvider:
//     meter.Int64Counter() → meterProvider.createCounter()
//
//  4. Метрики отправляются через Reader в MeterProvider:
//     Counter.Add() → Reader.collect() → Exporter.export() → OTLP Collector
//
// СХЕМА КОМПОНЕНТОВ:
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                     GLOBAL OTEL REGISTRY                           │
// │  otel.SetMeterProvider(provider) ← platform/metrics                │
// │  otel.Meter(name) → provider     ← assembly/metrics                     │
// └─────────────────────────────────────────────────────────────────────┘
//
//	↓
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                    METER PROVIDER (один)                           │
// │  ┌─────────────────────┐  ┌─────────────────────┐                  │
// │  │   Reader            │  │   Exporter          │                  │
// │  │ - Периодически      │  │ - Отправляет в      │                  │
// │  │   читает метрики    │  │   OTLP Collector    │                  │
// │  │ - Агрегирует        │  │ - Форматирует       │                  │
// │  │   данные            │  │   протокол          │                  │
// │  └─────────────────────┘  └─────────────────────┘                  │
// └─────────────────────────────────────────────────────────────────────┘
//
//	↓
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                     METERS (много)                                 │
// │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐     │
// │  │ assembly-service     │  │ database        │  │ http-client     │     │
// │  │ - RequestsTotal │  │ - Connections   │  │ - Requests      │     │
// │  │ - SightingsTotal│  │ - QueryDuration │  │ - Errors        │     │
// │  │ - AnalysisTime  │  │ - PoolSize      │  │ - Duration      │     │
// │  └─────────────────┘  └─────────────────┘  └─────────────────┘     │
// └─────────────────────────────────────────────────────────────────────┘
//
// ВАЖНЫЕ ПРИНЦИПЫ:
//
// 1. MeterProvider ОДИН - управляет инфраструктурой отправки метрик
// 2. Meter МНОГО - один на каждый логический компонент (сервис, библиотека)
// 3. Meter получает MeterProvider из глобального registry OpenTelemetry
// 4. Все метрики из всех Meter'ов отправляются через один MeterProvider
// 5. В Prometheus метрики группируются по label'у otel_scope_name
//
// Meter предоставляет методы для создания различных типов метрик:
// - Counter - монотонно возрастающий счетчик
// - UpDownCounter - счетчик, который может увеличиваться и уменьшаться
// - Histogram - распределение значений с bucketing
// - Gauge - моментальное значение (через UpDownCounter или Callback)
//
// Важно: meter должен быть создан один раз и переиспользоваться в рамках компонента
var meter = otel.Meter(serviceName)

// =============================================================================
// ТИПЫ МЕТРИК В OPENTELEMETRY
// =============================================================================
//
// 1. COUNTER (Счетчик) - metric.Int64Counter
//    - Монотонно возрастающее значение (только увеличивается)
//    - Используется для: количество запросов, ошибок, событий
//    - Пример: общее количество HTTP запросов
//    - Методы: Add() - добавить положительное значение
//
// 2. UPDOWNCOUNTER (Двунаправленный счетчик) - metric.Int64UpDownCounter
//    - Может увеличиваться и уменьшаться
//    - Используется для: активные соединения, размер очереди, память
//    - Пример: количество активных gRPC соединений
//    - Методы: Add() - добавить (может быть отрицательным)
//
// 3. HISTOGRAM (Гистограмма) - metric.Float64Histogram
//    - Распределение наблюдений в bucket'ах
//    - Автоматически создает метрики: _count, _sum, _bucket
//    - Используется для: время ответа, размер запроса, задержки
//    - Пример: время выполнения HTTP запроса
//    - Методы: Record() - записать наблюдение
//
// 4. GAUGE (Датчик) - НЕТ отдельного типа в OpenTelemetry!
//    - В OpenTelemetry нет прямого аналога Prometheus Gauge
//    - Для gauge-подобных метрик используются:
//      а) UpDownCounter - когда значение контролируется приложением
//      б) Асинхронные Observable - когда значение нужно читать по требованию
//    - Примеры: температура CPU, использование памяти, размер кэша
//    - Для простых случаев используйте UpDownCounter как gauge

var (
	// requestCounter - счетчик входящих запросов
	// Counter - монотонно возрастающее значение (никогда не уменьшается)
	// Используется для подсчета событий: запросы, ошибки, завершенные операции
	requestCounter metric.Int64Counter

	// responseCounter - счетчик исходящих ответов с атрибутами (статус, метод)
	// Позволяет группировать метрики по различным измерениям
	responseCounter metric.Int64Counter

	// AssembliesTotal - COUNTER для подсчета созданных заказов
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: бизнес-метрика для отслеживания количества новых заказов
	// Лейблы: нет (простой счетчик без группировки)
	AssembliesTotal metric.Int64Counter

	// AssembledDuration - HISTOGRAM для измерения времени выполнения сборки заказа
	// Тип: Float64Histogram (распределение значений)
	// Использование: отслеживание производительности алгоритма анализа
	// Автоматически создает метрики: _count, _sum, _bucket для percentile
	AssembledDuration metric.Float64Histogram
)

// InitMetrics инициализирует все метрики Assembly сервиса
// Должна быть вызвана один раз при старте приложения после инициализации OpenTelemetry провайдера
func InitMetrics() error {
	var err error

	// Создание счетчика запросов
	// Имя метрики следует соглашениям: namespace_protocol_app_metric_unit
	// WithDescription добавляет описание для документации и UI
	requestCounter, err = meter.Int64Counter(
		namespace+"_grpc_"+serviceName+"_requests_total",
		metric.WithDescription("Количество запросов к серверу"),
	)
	if err != nil {
		return err
	}

	// Создание счетчика ответов с поддержкой атрибутов
	// Атрибуты (labels в Prometheus) позволяют группировать метрики
	// Например: status="success", method="GetNote"
	responseCounter, err = meter.Int64Counter(
		namespace+"_grpc_"+serviceName+"_responses_total",
		metric.WithDescription("Количество ответов от сервера"),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик собранных заказов
	AssembliesTotal, err = meter.Int64Counter(
		namespace+"_"+serviceName+"_assembled_total",
		metric.WithDescription("Total number of Assembly orders assembled"),
	)
	if err != nil {
		return err
	}

	// Создаем гистограмму времени сборки заказа с правильными bucket'ами
	AssembledDuration, err = meter.Float64Histogram(
		namespace+"_"+serviceName+"_assembled_duration_seconds",
		metric.WithDescription("Duration of Assembly orders assembled"),
		metric.WithUnit("s"), // Указываем единицу измерения - секунды
		// Bucket'ы для бизнес-операций (анализ может занимать больше времени чем gRPC)
		metric.WithExplicitBucketBoundaries(
			0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0, 10.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}

// IncRequestCounter увеличивает счетчик входящих запросов на 1
// Вызывается при получении каждого нового запроса
func IncRequestCounter(ctx context.Context) {
	requestCounter.Add(ctx, 1)
}

// IncResponseCounter увеличивает счетчик ответов с указанными атрибутами
// status - статус ответа (success, error, timeout и т.д.)
// method - имя вызываемого метода gRPC
// Атрибуты позволяют создавать срезы метрик для анализа
func IncResponseCounter(ctx context.Context, status, method string) {
	responseCounter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("status", status),
			attribute.String("method", method),
		),
	)
}
