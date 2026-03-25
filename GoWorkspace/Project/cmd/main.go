package main

import (
	"Goworkspace/api/proto"
	"Goworkspace/internal/service"
	"Goworkspace/internal/storage"
	grpcServer "Goworkspace/internal/transport/grpc"
	transport "Goworkspace/internal/transport/http"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	st := storage.NewMemoryStorage()
	service := service.NewService(st)

	// --- HTTP Server ---
	httpRouter := transport.NewRouter(service)
	httpsrv := &http.Server{
		Addr:         ":8080",
		Handler:      httpRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("[INFO]: HTTP server started on :8080")
		if err := httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR]: HTTP listen error: %v", err)
		}
	}()

	// --- gRPC Server ---
	grpcSrv := grpcServer.NewTaskServer(service)
	grpcServer := grpc.NewServer()
	proto.RegisterTaskServiceServer(grpcServer, grpcSrv)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("[ERROR]: gRPC listen error: %v", err)
	}

	go func() {
		log.Println("[INFO]: gRPC server started on :50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("[ERROR]: gRPC serve error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("[INFO]: shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpsrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR]: HTTP graceful shutdown failed: %v", err)
	}

	grpcServer.GracefulStop()

	log.Println("[INFO]: servers stopped")
}
