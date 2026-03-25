package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// --- Создаём ключ для контекста ---
type contextKey string

// --- Делаем его константой чтобы нельзя было изменить ---
const requestIDKey contextKey = "request-id"

// --- Пишем функцию по созданию ID для запросов ---
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var requestID string

		// --- Вытаскиваем requestID из метадаты ---
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get("x-request-id"); len(ids) > 0 {
				requestID = ids[0]
			}
		}

		// --- Если requestID пустой, то генерируем новый ---
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// --- Кладём новый requestID в context ---
		ctx = context.WithValue(ctx, requestIDKey, requestID)
		return handler(ctx, req)
	}
}

// --- Пишем функцию для получения ID из context ---
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
