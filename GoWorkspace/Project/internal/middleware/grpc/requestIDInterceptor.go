package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// contextKey — приватный тип для ключей context, чтобы избежать конфликтов
type contextKey string

// requestIDKey — ключ для хранения request ID в context
const requestIDKey contextKey = "request-id"

// RequestIDInterceptor извлекает X-Request-ID из входящего запроса, генерирует новый, если отсутствует,
// и кладёт его в context для дальнейшего использования в сервисе
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var requestID string

		// Достаём X-Request-ID из метадаты gRPC; если нет — генерируем новый UUID
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get("x-request-id"); len(ids) > 0 {
				requestID = ids[0]
			}
		}
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Создаём новый context с requestID и вызываем handler с этим context
		ctx = context.WithValue(ctx, requestIDKey, requestID)
		return handler(ctx, req)
	}
}

// GetRequestID возвращает request ID из context или пустую строку
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
