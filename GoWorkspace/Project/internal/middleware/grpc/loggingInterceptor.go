package grpc

import (
	"Goworkspace/internal/logging"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// --- Пишем функцию по созданию логирования ---
func LoggingInterceptor(baseLogger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// --- Получем requestID из контекста и проверяем на пустоту ---
		requestID := GetRequestID(ctx)
		if requestID == "" {
			requestID = "unknown"
		}

		// --- Добавляем requestID и метод в логи ---
		logger := baseLogger.With(
			"request_id", requestID,
			"method", info.FullMethod,
		)

		// --- Создаём новый контекст с логированием + request id ---
		ctx = logging.WithLogger(ctx, logger)

		// --- Стартуем таймер запроса ---
		start := time.Now()
		logger.Info("gRPC: Request started")

		// --- Запускаем handler ---
		resp, err := handler(ctx, req)

		// --- Закрываем таймер и записываем длительность обработки ---
		duration := time.Since(start)

		// --- Проверяем на ошибки ---
		if err != nil {
			logger.Error("gRPC: request error",
				"duration", duration,
				"code", status.Code(err),
				"error", err,
			)
			return resp, err
		}

		logger.Info("gRPC: request success",
			"duration", duration,
		)
		return resp, nil
	}
}
