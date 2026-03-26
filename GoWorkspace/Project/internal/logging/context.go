package logging

import (
	"context"
	"log/slog"
)

// loggerKeyType — приватный тип ключа для context, чтобы избежать конфликтов
type loggerKeyType struct{}

// loggerKey — ключ для хранения logger в context
var loggerKey = loggerKeyType{}

// WithLogger возвращает новый context с добавленным slog.Logger
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext извлекает slog.Logger из context; если его нет, возвращает slog.Default()
func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
