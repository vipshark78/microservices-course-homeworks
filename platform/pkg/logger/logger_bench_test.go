package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func init() {
	InitForBenchmark()
}

func BenchmarkGlobalLogger(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(ctx, "test message")
	}
}

func BenchmarkWithLogger(b *testing.B) {
	log := With(zap.String("static_field", "static_value"))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(ctx, "test message")
	}
}

func BenchmarkChainLogger(b *testing.B) {
	type ContextValueKey string
	const traceKey ContextValueKey = "trace-123"
	const userKey ContextValueKey = "user-456"
	ctx := context.WithValue(context.Background(), traceIDKey, traceKey)
	ctx = context.WithValue(ctx, userIDKey, userKey)
	log := With(zap.String("static_field", "static_value"))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info(ctx, "test message")
	}
}
