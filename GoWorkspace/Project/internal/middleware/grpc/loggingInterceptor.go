package inteceptor

import (
	"Goworkspace/internal/logging"
	"context"
	"log"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor(baseLogger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		requestID := GetRequestID(ctx)
		reqLogger := baseLogger.With(
			"request-id", requestID,
			"method", info.FullMethod,
		)

		ctx = logging.WithLogger(ctx, reqLogger)

		log.Printf("[INFO]: gRPC: request started id=%s, method=%s", requestID, info.FullMethod)

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			log.Printf("[ERROR]: gRPC: id=%s, method=%s, duration=%s, error=%v", requestID, info.FullMethod, duration, err)
		}

		log.Printf("[INFO]: gRPC: id=%s, method=%s, duration=%s", requestID, info.FullMethod, duration)

		return resp, err
	}
}
