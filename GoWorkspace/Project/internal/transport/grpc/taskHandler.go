package grpcTransport

import (
	"Goworkspace/api/proto"
	"Goworkspace/internal/domain"
	"context"
)

type TaskServer struct {
	proto.UnimplementedTaskServiceServer
	service domain.TaskService
}

func NewTaskServer(service domain.TaskService) *TaskServer {
	return &TaskServer{service: service}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *proto.CreateTaskRequest) (*proto.CreateTaskResponse, error) {
	item, err := s.service.Create(ctx, req.Name)
	if err != nil {
		return nil, HelperErrorGRPC(ctx, err)
	}

	return &proto.CreateTaskResponse{
		Task: itemToProto(item),
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.GetTaskResponse, error) {
	item, err := s.service.Get(ctx, int(req.Id))
	if err != nil {
		return nil, HelperErrorGRPC(ctx, err, req.Id)
	}

	return &proto.GetTaskResponse{
		Task: itemToProto(item),
	}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *proto.DeleteTaskRequest) (*proto.DeleteTaskResponse, error) {
	err := s.service.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, HelperErrorGRPC(ctx, err, req.Id)
	}

	return &proto.DeleteTaskResponse{
		Id: req.Id,
	}, nil
}
