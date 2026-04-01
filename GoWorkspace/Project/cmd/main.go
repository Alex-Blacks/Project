package main

import (
	"Goworkspace/api/proto"
	"Goworkspace/internal/logging"
	grpcMiddleware "Goworkspace/internal/middleware/grpc"
	"Goworkspace/internal/service"
	"Goworkspace/internal/service/auth"
	"Goworkspace/internal/storage"
	grpcTransport "Goworkspace/internal/transport/grpc"
	transport "Goworkspace/internal/transport/http"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := logging.NewLogger()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logger.Error("JWT_SECRET is not set")
		return
	}
	jwt := auth.NewJWTValidation([]byte(secret))
	authService := auth.NewAuthService(jwt)
	public := map[string]struct{}{
		"/proto.AuthService/Login": {},
	}
	st := storage.NewMemoryStorage()
	svc := service.NewService(st)

	// --- HTTP Server ---
	httpRouter := transport.NewRouter(svc)
	httpsrv := &http.Server{
		Addr:         ":8080",
		Handler:      httpRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("HTTP: server started on :8080")
		if err := httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP: listen error", "error", err)
			return
		}
	}()

	// --- gRPC Server ---
	grpcSrv := grpcTransport.NewTaskServer(svc)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcMiddleware.RequestIDInterceptor(),
			grpcMiddleware.TimeoutInterceptor(60*time.Second),
			grpcMiddleware.LoggingInterceptor(logger),
			grpcMiddleware.RecoveryInterceptor(),
			grpcMiddleware.NewAuthInterceptor(authService, public).Unary(),
		))
	reflection.Register(grpcServer)
	authHandler := grpcTransport.NewAuthHandler(authService)
	proto.RegisterAuthServiceServer(grpcServer, authHandler)
	proto.RegisterTaskServiceServer(grpcServer, grpcSrv)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("gRPC: listen error", "error", err)
		return
	}

	go func() {
		logger.Info("gRPC: server started on :50051")
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC: serve error", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpsrv.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP: graceful shutdown failed", "error", err)
	}

	grpcServer.GracefulStop()

	logger.Info("servers stopped")
}
