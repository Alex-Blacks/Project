package grpcTransport

import (
	"Goworkspace/api/proto"
	"Goworkspace/internal/domain"
	"Goworkspace/internal/logging"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HelperErrorGRPC(ctx context.Context, err error, id ...int64) error {
	logger := logging.LoggerFromContext(ctx)
	if len(id) > 0 {
		logger.Error("gRPC:",
			"id", id[0],
			"error", err,
		)
	} else {
		logger.Error("gRPC:",
			"error", err,
		)
	}
	return MapDomainErrorToCodes(err)
}

func MapDomainErrorToCodes(err error) error {
	switch err {
	case domain.ErrInvalidValue,
		domain.ErrEmptyName:
		return status.Error(codes.InvalidArgument, err.Error())
	case domain.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func itemToProto(item domain.Item) *proto.Task {
	return &proto.Task{
		Id:   int64(item.ID),
		Name: item.Name,
	}
}
