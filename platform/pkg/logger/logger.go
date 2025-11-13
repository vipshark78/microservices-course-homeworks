package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Key             string
	ContextValueKey string
)

const (
	traceIDKey ContextValueKey = "trace_id"
	userIDKey  ContextValueKey = "user_id"
)

var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
)

// logger обёртка над zap-логгером
type logger struct {
	zapLogger *zap.Logger
}

func Init(level string, asJSON bool) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLvl(level))
		encoderCfg := buildProductionEncoderConfig()
		var encoder zapcore.Encoder
		if asJSON {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}
		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			dynamicLevel,
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})
	return nil
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
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}
	return nil
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
