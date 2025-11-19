package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

// shutdownTimeout по умолчанию 10 секунд
const shutdownTimeout = 10 * time.Second

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	done   chan struct{}
	funcs  []func(ctx context.Context) error
	logger Logger
}

// globalCloser глобальный экземпляр Closer
var globalCloser = NewWithLogger(&logger.NoopLogger{})

// AddNamed добавляет функцию закрытия с именем зависимости для логирования в глобальный closer
func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

// Add добавляет функции закрытия в глобальный closer
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// CloseAll инициирует процесс закрытия всех зарегистрированных функций глобального closer'а
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

// SetLogger позволяет установить кастомный логгер для глобального closer'a
func SetLogger(l Logger) {
	globalCloser.SetLogger(l)
}

// Configure настраивает глобальный closer для обработки системных сигналов
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

// New создает новый экземпляр closer с дефолтным логгером log.Default()
func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

// NewWithLogger создает новый экземпляр closer с заданным логгером
// Если переданы сигналы, Closer начнет их слушать и вызывать CloseAll при получении любого из них.
func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// SetLogger устанавливает логгер для Closer'а
func (c *Closer) SetLogger(logger Logger) {
	c.logger = logger
}

// handleSignals обрабатывает системные сигналы и вызывает CloseAll с fresh shutdown context
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "🛑 Получен системный сигнал, начинаем graceful shutdown...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "❌ Ошибка при закрытии ресурсов: %v", zap.Error(err))
		}

	case <-c.done:
		// CloseAll уже был вызван вручную, просто выходим
	}
}

// AddNamed добавляет функцию закрытия с именем зависимости для логирования
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("🧩 Закрываем %s...", name))

		err := f(ctx)

		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("❌ Ошибка при закрытии %s: %v (заняло %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("✅ %s успешно закрыт за %s", name, duration))
		}
		return err
	})
}

// Add добавляет одну или несколько функций закрытия в Closer
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "ℹ️ Нет функций для закрытия.")
			return
		}

		c.logger.Info(ctx, "🚦 Начинаем процесс graceful shutdown...")
		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		// Выполняем в обратном порядке добавления
		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				// Защита от паники
				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.logger.Error(ctx, "⚠️ Panic в функции закрытия", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		// Закрываем канал ошибок, когда все функции завершатся
		go func() {
			wg.Wait()
			close(errCh)
		}()

		// Читаем ошибки или отмену контекста
		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "⚠️ Контекст отменён во время закрытия", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "✅ Все ресурсы успешно закрыты")
					return
				}
				c.logger.Error(ctx, "❌ Ошибка при закрытии", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
