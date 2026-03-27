package grpc

import (
	"Goworkspace/internal/logging"
	"Goworkspace/internal/service"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ключ для хранения userID в context
// используем свой тип (contextKey), чтобы избежать коллизий
const userIDKey contextKey = "user_id"

// interceptor хранит зависимость AuthService
// через него происходит проверка токена
type AuthInterceptor struct {
	auth service.AuthService
}

// Unary возвращает gRPC interceptor
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		// достаём logger из context
		logger := logging.LoggerFromContext(ctx)

		var token string

		// извлекаем metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// если metadata нет — сразу ошибка
			logger.Error("gRPC: missing metadata")
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// получаем header "authorization"
		values := md.Get("authorization")

		if len(values) == 0 {
			// header отсутствует
			logger.Error("gRPC: missing authorization header")
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		// ожидаем формат: "Bearer <token>"
		parts := strings.SplitN(values[0], " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			// неправильный формат header
			logger.Error("gRPC: invalid authorization header")
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
		}

		// извлекаем сам токен
		token = parts[1]

		// валидируем токен через сервис (JWT внутри)
		userID, err := i.auth.Validate(token)
		if err != nil || userID == "" {
			// токен невалидный или не удалось извлечь userID
			logger.Error("gRPC: invalid token")
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// кладём userID в context, чтобы handler мог его использовать
		ctx = context.WithValue(ctx, userIDKey, userID)

		// передаём управление дальше (в следующий interceptor или handler)
		return handler(ctx, req)
	}
}

// GetUserID возвращает user ID из context или пустую строку
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}
