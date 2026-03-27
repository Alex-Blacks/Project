package service

import (
	"Goworkspace/internal/domain"
	"context"
	"errors"
)

type Service struct {
	storage domain.Storage
}

func NewService(st domain.Storage) *Service {
	return &Service{storage: st}
}

var _ domain.TaskService = (*Service)(nil)

func (s *Service) Create(ctx context.Context, name string) (domain.Item, error) {
	if name == "" {
		return domain.Item{}, domain.ErrEmptyName
	}

	itemName := domain.Item{Name: name}
	item, err := s.storage.CreateItem(ctx, itemName)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return domain.Item{}, err
		}
		return domain.Item{}, domain.ErrInternal
	}
	return item, nil
}

func (s *Service) Get(ctx context.Context, id int) (domain.Item, error) {
	if id < 1 {
		return domain.Item{}, domain.ErrInvalidValue
	}

	item, err := s.storage.GetItem(ctx, id)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return domain.Item{}, err
		}
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Item{}, domain.ErrNotFound
		}
		return domain.Item{}, domain.ErrInternal
	}

	return item, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return domain.ErrInvalidValue
	}

	if err := s.storage.DeleteItem(ctx, id); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return domain.ErrInternal
	}

	return nil
}
