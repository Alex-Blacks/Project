package storage

import (
	"context"
	"sync"

	"Goworkspace/Project/domain"
)

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[int]domain.Item
	next int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int]domain.Item),
		next: 1,
	}
}

func (s *MemoryStorage) CreateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	select {
	case <-ctx.Done():
		return domain.Item{}, ctx.Err()
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		item.ID = s.next
		s.next++
		s.data[item.ID] = item

		return item, nil
	}

}

func (s *MemoryStorage) GetItem(ctx context.Context, id int) (domain.Item, error) {
	select {
	case <-ctx.Done():
		return domain.Item{}, ctx.Err()
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()

		item, ok := s.data[id]

		if !ok {
			return domain.Item{}, domain.ErrNotFound
		}

		return item, nil
	}
}

func (s *MemoryStorage) DeleteItem(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.mu.Lock()
		defer s.mu.Unlock()

		if _, ok := s.data[id]; !ok {
			return domain.ErrNotFound
		}

		delete(s.data, id)

		return nil
	}

}
