package domain_test

import (
	"context"
	"errors"
	"testing"

	"Goworkspace/Project/domain"
)

type MockStorage struct {
	storageCalled bool
	forcedError   error
}

func (m *MockStorage) CreateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	m.storageCalled = true
	return item, m.forcedError
}
func (m *MockStorage) GetItem(ctx context.Context, id int) (domain.Item, error) {
	m.storageCalled = true
	return domain.Item{ID: id}, m.forcedError
}
func (m *MockStorage) DeleteItem(ctx context.Context, id int) error {
	m.storageCalled = true
	return m.forcedError
}

func TestService_Create(t *testing.T) {
	t.Run("Empty name returns ErrEmptyName", func(t *testing.T) {
		mock := &MockStorage{}
		svc := domain.NewService(mock)

		_, err := svc.Create(context.Background(), "")
		if !errors.Is(err, domain.ErrEmptyName) {
			t.Fatalf("expected ErrEmptyName, got %v", err)
		}
		if mock.storageCalled {
			t.Fatal("storage should not be called for empty name")
		}
	})

	t.Run("Success returns item", func(t *testing.T) {
		mock := &MockStorage{}
		svc := domain.NewService(mock)

		item, err := svc.Create(context.Background(), "Alex")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.Name != "Alex" {
			t.Fatalf("expected name Alex, got %v", item.Name)
		}
		if !mock.storageCalled {
			t.Fatal("storage should be called")
		}
	})

	t.Run("Storage error returns ErrInternal", func(t *testing.T) {
		mock := &MockStorage{forcedError: errors.New("db fail")}
		svc := domain.NewService(mock)

		_, err := svc.Create(context.Background(), "Alex")
		if !errors.Is(err, domain.ErrInternal) {
			t.Fatalf("expected ErrInternal, got %v", err)
		}
	})

	t.Run("Context canceled returns context.Canceled", func(t *testing.T) {
		mock := &MockStorage{forcedError: context.Canceled}
		svc := domain.NewService(mock)

		_, err := svc.Create(context.Background(), "Alex")
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Context timeout returns context.DeadlineExceeded", func(t *testing.T) {
		mock := &MockStorage{forcedError: context.DeadlineExceeded}
		svc := domain.NewService(mock)

		_, err := svc.Create(context.Background(), "Alex")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.DeadlineExceeded, got %v", err)
		}
	})
}

func TestService_Get(t *testing.T) {
	t.Run("Zero ID returns ErrInvalidValue", func(t *testing.T) {
		mock := &MockStorage{}
		service := domain.NewService(mock)

		_, err := service.Get(context.Background(), 0)
		if !errors.Is(err, domain.ErrInvalidValue) {
			t.Fatalf("expected ErrInvalidValue, got: %v", err)
		}
	})

	t.Run("Negative ID returns ErrInvalidValue", func(t *testing.T) {
		mock := &MockStorage{}
		service := domain.NewService(mock)

		_, err := service.Get(context.Background(), -1)
		if !errors.Is(err, domain.ErrInvalidValue) {
			t.Fatalf("expected ErrInvalidValue, got: %v", err)
		}
	})

	t.Run("Success returns item", func(t *testing.T) {
		mock := &MockStorage{}
		service := domain.NewService(mock)

		item, err := service.Get(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item.ID != 1 {
			t.Fatalf("expected item id:1, got: %v", err)
		}
	})

	t.Run("Storage error returns ErrNotFound", func(t *testing.T) {
		mock := &MockStorage{forcedError: domain.ErrNotFound}
		service := domain.NewService(mock)

		_, err := service.Get(context.Background(), 1)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got: %v", err)
		}
	})

	t.Run("Storage error returns ErrInternal", func(t *testing.T) {
		mock := &MockStorage{forcedError: errors.New("DB error")}
		service := domain.NewService(mock)

		_, err := service.Get(context.Background(), 1)
		if !errors.Is(err, domain.ErrInternal) {
			t.Fatalf("expected ErrInternal, got: %v", err)
		}
	})

	t.Run("Context canceled returns context.Canceled", func(t *testing.T) {
		mock := &MockStorage{forcedError: context.Canceled}
		svc := domain.NewService(mock)

		_, err := svc.Get(context.Background(), 1)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Context timeout returns context.DeadlineExceeded", func(t *testing.T) {
		mock := &MockStorage{forcedError: context.DeadlineExceeded}
		svc := domain.NewService(mock)

		_, err := svc.Get(context.Background(), 1)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.DeadlineExceeded, got %v", err)
		}
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("Zero ID returns ErrInvalidValue", func(t *testing.T) {
		mock := &MockStorage{}
		service := domain.NewService(mock)

		err := service.Delete(context.Background(), 0)
		if !errors.Is(err, domain.ErrInvalidValue) {
			t.Fatalf("expected ErrInvalidValue, got: %v", err)
		}
	})

	t.Run("Negative ID returns ErrInvalidValue", func(t *testing.T) {
		mock := &MockStorage{}
		service := domain.NewService(mock)

		err := service.Delete(context.Background(), -1)
		if !errors.Is(err, domain.ErrInvalidValue) {
			t.Fatalf("expected ErrInvalidValue, got: %v", err)
		}
	})

	t.Run("Storage error returns ErrNotFound", func(t *testing.T) {
		mock := &MockStorage{forcedError: domain.ErrNotFound}
		service := domain.NewService(mock)

		err := service.Delete(context.Background(), 1)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got: %v", err)
		}
	})

	t.Run("Storage error returns ErrInternal", func(t *testing.T) {
		mock := &MockStorage{forcedError: errors.New("DB error")}
		service := domain.NewService(mock)

		err := service.Delete(context.Background(), 1)
		if !errors.Is(err, domain.ErrInternal) {
			t.Fatalf("expected ErrInternal, got: %v", err)
		}
	})

	t.Run("Context canceled returns context.Canceled", func(t *testing.T) {
		mock := &MockStorage{forcedError: context.Canceled}
		service := domain.NewService(mock)

		err := service.Delete(context.Background(), 1)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got: %v", err)
		}
	})

	t.Run("Context timeout returns context.DeadlineExceeded", func(t *testing.T) {
		mock := &MockStorage{forcedError: context.DeadlineExceeded}
		service := domain.NewService(mock)

		err := service.Delete(context.Background(), 1)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.DeadlineExceeded, got: %v", err)
		}
	})
}
