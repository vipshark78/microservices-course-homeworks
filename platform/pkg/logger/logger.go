package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Key             string
	ContextValueKey string
)

const (
	traceIDKey         ContextValueKey = "trace_id"
	userIDKey          ContextValueKey = "user_id"
	serviceEnvironment                 = "dev"           // окружение для фильтрации логов
	shutdownTimeout                    = 2 * time.Second // таймаут для graceful shutdown OTLP provider
)

var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
	otelProvider *otelLogSdk.LoggerProvider
)

// logger обёртка над zap-логгером
type logger struct {
	zapLogger *zap.Logger
}

// Init инициализирует глобальный логгер с Tee архитектурой.
// Поддерживает одновременную запись в stdout и OTLP коллектор.
//
// Параметры:
//   - logLevel: уровень логирования ("debug", "info", "warn", "error")
//   - asJSON: формат вывода (true - JSON, false - консольный)
//   - enableOTLP: включение отправки в OpenTelemetry коллектор
func Init(ctx context.Context, logLevel string, asJSON, enableOTLP bool, otlpEndpoint, serviceName string) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLvl(logLevel))
		cores := buildCores(ctx, asJSON, enableOTLP, otlpEndpoint, serviceName)
		zapLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})

	if globalLogger == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

// buildCores создает слайс cores для zapcore.Tee.
// Всегда включает stdout core, опционально добавляет OTLP core.
func buildCores(ctx context.Context, asJSON, enableOTLP bool, otlpEndpoint, serviceName string) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(asJSON),
	}

	if enableOTLP {
		if otlpCore := createOTLPCore(ctx, otlpEndpoint, serviceName); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr.
// Поддерживает JSON и консольный формат вывода.
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildProductionEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), dynamicLevel)
}

// createOTLPCore создает core для отправки в OpenTelemetry коллектор.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCore(ctx context.Context, otlpEndpoint, serviceName string) *SimpleOTLPCore {
	otlpLogger, err := createOTLPLogger(ctx, otlpEndpoint, serviceName)
	if err != nil {
		// Логирование ошибки невозможно, так как логгер еще не инициализирован
		return nil
	}

	// Прямо передаём OTLP-логгер в core. Буферизацию делает OTLP SDK (BatchProcessor).
	return NewSimpleOTLPCore(otlpLogger, dynamicLevel)
}

// createOTLPLogger создает OTLP логгер с настроенным экспортером и ресурсами.
// Использует BatchProcessor для эффективной отправки логов.
func createOTLPLogger(ctx context.Context, endpoint, serviceName string) (otelLog.Logger, error) {
	exporter, err := createOTLPExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createResource(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)
	otelProvider = provider // сохраняем для shutdown

	return provider.Logger("app"), nil
}

// createOTLPExporter создает gRPC экспортер для OTLP коллектора
func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // для разработки, в продакшене следует использовать TLS
	)
}

// createResource создает метаданные сервиса для телеметрии
func createResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", serviceEnvironment),
		),
	)
}

func InitForBenchmark() {
	core := zapcore.NewNopCore()

	globalLogger = &logger{
		zapLogger: zap.New(core),
	}
}

// SetLevel устанавливает уровень логирования
func SetLevel(lvl string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}
	dynamicLevel.SetLevel(parseLvl(lvl))
}

// Logger возвращает глобальный логгер
func Logger() *logger {
	return globalLogger
}

// SetNopLogger устанавливает глобальный логгер в nop режиме
func SetNopLogger() {
	globalLogger = &logger{
		zapLogger: zap.NewNop(),
	}
}

// Sync сбрасывает буферы логгера
func Sync(_ context.Context) error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}
	return nil
}

// Close корректно завершает работу логгера.
// Останавливает OTLP provider с таймаутом для отправки оставшихся логов.
func Close(ctx context.Context) error {
	if otelProvider != nil {
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		err := otelProvider.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("Error shutting down OTLP provider:%w", err)
		}
	}

	return nil
}

func CloseAndSync(ctx context.Context) error {
	err := Close(ctx)
	if errSync := Sync(ctx); errSync != nil {
		if err != nil {
			return fmt.Errorf("close error: %w, sync error: %w", err, errSync)
		}
		return errSync
	}
	return err
}

// With создает новый логгер с переданными полями
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}
	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext создает логгер с контекстом
func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",                           // время
		LevelKey:       "level",                        // уровень логирования
		NameKey:        "logger",                       // имя логгера
		CallerKey:      "caller",                       // место вызова
		MessageKey:     "msg",                          // сообщение
		StacktraceKey:  "stacktrace",                   // стектрейс
		LineEnding:     zapcore.DefaultLineEnding,      // разделитель строк в консоли
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // уровень логгирования капсом
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // формат времени ISO8601
		EncodeDuration: zapcore.SecondsDurationEncoder, // длительность в секундах
		EncodeCaller:   zapcore.ShortCallerEncoder,     // краткий путь до файла и номер строки
		EncodeName:     zapcore.FullNameEncoder,        // полное имя логгера
	}
}

func parseLvl(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
