package inteceptor

import (
	"context"
	"log"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func loggingInterception(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		requestID := GetRequestID(ctx)
		reqLogger := logger.With("request-id", requestID)

		ctx = context.WithValue(ctx, loggerKey, reqLogger)
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
