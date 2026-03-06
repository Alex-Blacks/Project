package grpc

import (
	"Goworkspace/Project/api/proto"
	"Goworkspace/Project/domain"
	"context"
)

type TaskServer struct {
	proto.UnimplementedTaskServiceServer
	service *domain.Service
}

func NewTaskServer(service *domain.Service) *TaskServer {
	return &TaskServer{service: service}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *proto.CreateTaskRequest) (*proto.Task, error) {
	item, err := s.service.Create(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &proto.Task{
		Id:   int64(item.ID),
		Name: item.Name,
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.Task, error) {
	item, err := s.service.Get(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &proto.Task{
		Id:   int64(item.ID),
		Name: item.Name,
	}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *proto.DeleteTaskRequest) (*proto.DeleteTaskResponse, error) {
	err := s.service.Delete(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &proto.DeleteTaskResponse{
		Id: req.Id,
	}, nil
}
