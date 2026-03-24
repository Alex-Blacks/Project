package grpcServer_test

import (
	"Goworkspace/api/proto"
	"Goworkspace/internal/service"
	"Goworkspace/internal/storage"
	grpcPkg "Goworkspace/internal/transport/grpc"

	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func setupServer() *grpc.ClientConn {
	lis := bufconn.Listen(bufSize)

	st := storage.NewMemoryStorage()
	service := service.NewService(st)

	server := grpc.NewServer()
	taskServer := grpcPkg.NewTaskServer(service)
	proto.RegisterTaskServiceServer(server, taskServer)

	go func() {
		server.Serve(lis)
	}()

	conn, _ := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	return conn
}

func TestCreateTask(t *testing.T) {
	conn := setupServer()
	defer conn.Close()
	client := proto.NewTaskServiceClient(conn)

	t.Run("Create: success create", func(t *testing.T) {
		resp, err := client.CreateTask(context.Background(), &proto.CreateTaskRequest{
			Name: "Alex",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Task.Name != "Alex" && resp.Task.Id != 1 {
			t.Fatalf("Expected name: Alex, id: 1, got name: %s, id: %d", resp.Task.Name, resp.Task.Id)
		}
	})
	t.Run("Create: empty name", func(t *testing.T) {
		_, err := client.CreateTask(context.Background(), &proto.CreateTaskRequest{
			Name: "",
		})
		if err == nil {
			t.Fatal("Expected error")
		}
		if status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
		}
	})
}

func TestGetTask(t *testing.T) {
	conn := setupServer()
	defer conn.Close()
	client := proto.NewTaskServiceClient(conn)

	t.Run("Get: success Get", func(t *testing.T) {
		CreateResp, err := client.CreateTask(context.Background(), &proto.CreateTaskRequest{Name: "Alex"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		GetResp, err := client.GetTask(context.Background(), &proto.GetTaskRequest{Id: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if GetResp.Task.Name != CreateResp.Task.Name && GetResp.Task.Id != CreateResp.Task.Id {
			t.Fatalf("Expected name: %s, id: %d, got name: %s, id: %d", CreateResp.Task.Name, CreateResp.Task.Id, GetResp.Task.Name, GetResp.Task.Id)
		}
	})
	t.Run("Get: Invalide value", func(t *testing.T) {
		_, err := client.GetTask(context.Background(), &proto.GetTaskRequest{Id: 0})
		if status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
		}
	})
	t.Run("Get: Not Found", func(t *testing.T) {
		_, err := client.GetTask(context.Background(), &proto.GetTaskRequest{Id: 2})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("expected NotFound, got %v", status.Code(err))
		}
	})
}
func TestDeleteTask(t *testing.T) {
	conn := setupServer()
	defer conn.Close()
	client := proto.NewTaskServiceClient(conn)

	t.Run("Delete: success delete", func(t *testing.T) {
		CreateResp, err := client.CreateTask(context.Background(), &proto.CreateTaskRequest{Name: "Alex"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		DeleteResp, err := client.DeleteTask(context.Background(), &proto.DeleteTaskRequest{Id: 1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if CreateResp.Task.Id != DeleteResp.Id {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("Delete: Invalide value", func(t *testing.T) {
		_, err := client.DeleteTask(context.Background(), &proto.DeleteTaskRequest{Id: 0})
		if status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
		}
	})
	t.Run("Delete: Not Found", func(t *testing.T) {
		_, err := client.DeleteTask(context.Background(), &proto.DeleteTaskRequest{Id: 2})
		if status.Code(err) != codes.NotFound {
			t.Fatalf("expected NotFound, got %v", status.Code(err))
		}
	})
}
