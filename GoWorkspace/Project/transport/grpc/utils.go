package grpc

import (
	"Goworkspace/Project/domain"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HelperErrorGRPC(err error, id ...int64) error {
	if len(id) > 0 {
		log.Printf("[ERROR]: gRPC: id=%d: %v", id[0], err)
	} else {
		log.Printf("[ERROR]: gRPC: %v", err)
	}
	return MapDomainErrorToCodes(err)
}

func MapDomainErrorToCodes(err error) error {
	switch err {
	case domain.ErrInvalidValue:
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
