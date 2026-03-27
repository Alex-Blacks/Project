package grpc

import (
	"Goworkspace/internal/logging"
	"context"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryInterceptor перехватывает panic в gRPC handler'ах,
// логирует детали (включая stack trace) и возвращает клиенту codes.Internal
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		// defer выполняется после handler и позволяет перехватить panic
		defer func() {
			if rec := recover(); rec != nil {
				// Получаем logger из context (уже содержит request_id и method)
				logger := logging.LoggerFromContext(ctx)

				// Логируем panic и stack trace для диагностики
				logger.Error("gRPC: panic recovered",
					"panic", rec,
					"stack", string(debug.Stack()),
				)

				// Возвращаем клиенту безопасную ошибку без внутренних деталей
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		// Вызываем следующий handler
		return handler(ctx, req)
	}
}
