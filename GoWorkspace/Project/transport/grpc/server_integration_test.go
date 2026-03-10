package grpcPkg_test

import (
	"Goworkspace/Project/api/proto"
	"Goworkspace/Project/domain"
	"Goworkspace/Project/storage"
	grpcPkg "Goworkspace/Project/transport/grpc"
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
	service := domain.NewService(st)

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
		if resp.Name != "Alex" && resp.Id != 1 {
			t.Fatalf("Expected name: Alex, id: 1, got name: %s, id: %d", resp.Name, resp.Id)
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

		if GetResp.Name != CreateResp.Name && GetResp.Id != CreateResp.Id {
			t.Fatalf("Expected name: %s, id: %d, got name: %s, id: %d", CreateResp.Name, CreateResp.Id, GetResp.Name, GetResp.Id)
		}
	})
	t.Run("Get: Invalide ID", func(t *testing.T) {
		_, err := client.GetTask(context.Background(), &proto.GetTaskRequest{Id: 0})
		if status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
		}
	})
}

func TestInvalidValue(t *testing.T) {
	conn := setupServer()
	defer conn.Close()

	client := proto.NewTaskServiceClient(conn)

	_, err := client.GetTask(
		context.Background(),
		&proto.GetTaskRequest{
			Id: -2,
		},
	)
	if err == nil {
		t.Fatal("Expected error")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}
