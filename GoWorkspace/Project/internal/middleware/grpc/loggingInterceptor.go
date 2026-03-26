package grpc

import (
	"Goworkspace/internal/logging"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor логирует gRPC запросы: старт, успешное завершение, ошибки и длительность
func LoggingInterceptor(baseLogger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Получаем requestID из context; если нет, используем "unknown"
		requestID := GetRequestID(ctx)
		if requestID == "" {
			requestID = "unknown"
		}

		// Создаём logger с requestID и gRPC методом для структурированного логирования
		logger := baseLogger.With(
			"request_id", requestID,
			"method", info.FullMethod,
		)

		// Добавляем logger в context, чтобы сервисы могли его использовать
		ctx = logging.WithLogger(ctx, logger)

		// Засекаем время обработки и логируем старт запроса
		start := time.Now()
		logger.Info("gRPC: Request started")

		// Вызываем следующий handler с обновлённым context
		resp, err := handler(ctx, req)

		// Логируем длительность запроса; если есть ошибка — логируем с кодом ошибки, иначе — логируем успех
		duration := time.Since(start)
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
