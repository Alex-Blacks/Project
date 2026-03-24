package inteceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const requestIDKey contextKey = "request-id"

func GetRequestID(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func RequestIDInteceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var requestID string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if ids := md.Get("x-request-id"); len(ids) > 0 {
			requestID = ids[0]
		}
	}

	if requestID == "" {
		requestID = uuid.New().String()
	}
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	return handler(ctx, req)
}
