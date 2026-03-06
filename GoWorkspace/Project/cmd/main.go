package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Goworkspace/Project/api/proto"
	"Goworkspace/Project/domain"
	"Goworkspace/Project/storage"
	grpcPkg "Goworkspace/Project/transport/grpc"
	transport "Goworkspace/Project/transport/http"

	"google.golang.org/grpc"
)

func main() {
	st := storage.NewMemoryStorage()
	service := domain.NewService(st)

	// --- HTTP Server ---
	r := transport.NewRouter(service)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("[INFO]: http server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR]: listen error: %v", err)
		}
	}()

	// --- gRPC Server ---
	grpcSrv := grpcPkg.NewTaskServer(service)
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

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	<-stop

	log.Println("[INFO]: shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR]: graceful shutdown failed: %v", err)
	}

	log.Println("[INFO]: server stopped")
}
