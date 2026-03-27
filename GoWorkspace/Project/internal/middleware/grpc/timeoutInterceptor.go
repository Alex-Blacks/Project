package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutInterceptor создаёт interceptor, который проверяет deadline клиента и если он больше серверного то ограничивает время обработки запроса.
// Если handler не завершится за указанный timeout, context автоматически отменяется.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Проверяем что клиентеский deadline меньше серверного и запускаем с ним handler
		if deadline, ok := ctx.Deadline(); ok {
			if time.Until(deadline) <= timeout {
				return handler(ctx, req)
			}
		}

		// Добавляем свой таймаут в контекст, если он меньше серверного
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return handler(ctx, req)
	}
}
