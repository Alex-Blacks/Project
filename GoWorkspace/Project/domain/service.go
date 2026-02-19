package domain

import (
	"context"
	"errors"
)

type Service struct {
	storage Storage
}

func NewService(st Storage) *Service {
	return &Service{storage: st}
}

type Item struct {
	ID   int
	Name string
}

func (s *Service) Create(ctx context.Context, name string) (Item, error) {
	if name == "" {
		return Item{}, ErrEmptyName
	}

	itemName := Item{Name: name}
	item, err := s.storage.CreateItem(ctx, itemName)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return Item{}, err
		}
		return Item{}, ErrInternal
	}
	return item, nil
}

func (s *Service) Get(ctx context.Context, id int) (Item, error) {
	if id < 1 {
		return Item{}, ErrInvalidValue
	}

	item, err := s.storage.GetItem(ctx, id)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return Item{}, err
		}
		if errors.Is(err, ErrNotFound) {
			return Item{}, ErrNotFound
		}
		return Item{}, ErrInternal
	}

	return item, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return ErrInvalidValue
	}

	if err := s.storage.DeleteItem(ctx, id); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return ErrInternal
	}

	return nil
}
