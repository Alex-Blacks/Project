package grpcTransport

import (
	"Goworkspace/api/proto"
	"Goworkspace/internal/logging"
	"Goworkspace/internal/service/auth"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	proto.UnimplementedAuthServiceServer
	auth auth.AuthService
}

func NewAuthHandler(auth auth.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (a *AuthHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	logger := logging.LoggerFromContext(ctx)
	token, err := a.auth.Login(req.Login, req.Password)
	if err != nil {
		logger.Error("invalid credentials",
			"error", err,
		)
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials error:%v", err)
	}

	return &proto.LoginResponse{Token: token}, nil
}
