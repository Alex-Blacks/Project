package storage_test

import (
	"Goworkspace/Project/domain"
	"Goworkspace/Project/storage"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func CanceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return ctx
}
func TimeoutContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	<-ctx.Done()

	return ctx
}

func TestStorage_Create(t *testing.T) {

	t.Run("Success return item", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		item := domain.Item{Name: "Alex"}

		resItem, err := st.CreateItem(context.Background(), item)
		if resItem.ID != 1 || resItem.Name != "Alex" {
			t.Fatalf("expected item id=1, name: Alex; got: id: %v, name: %v", resItem.ID, resItem.Name)
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("Increment ID by 1", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		st.CreateItem(context.Background(), domain.Item{Name: "Alex"})

		resItem, err := st.CreateItem(context.Background(), domain.Item{Name: "Alice"})
		if resItem.ID != 2 || resItem.Name != "Alice" {
			t.Fatalf("expected item id= 2, name: Alice; got: id: %v, name: %v", resItem.ID, resItem.Name)
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("Data save in Memory Storage", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		st.CreateItem(context.Background(), domain.Item{Name: "Alex"})
		st.CreateItem(context.Background(), domain.Item{Name: "Alice"})

		resItem, err := st.GetItem(context.Background(), 2)

		if resItem.ID != 2 || resItem.Name != "Alice" {
			t.Fatalf("expected item id= 2, name: Alice; got: id: %v, name: %v", resItem.ID, resItem.Name)
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("Context Canceled returns context.Canceled", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		item := domain.Item{Name: "Deril"}

		_, err := st.CreateItem(CanceledContext(), item)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Context timeout returns context.DeadlineExceeded", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		item := domain.Item{Name: "Deril"}

		_, err := st.CreateItem(TimeoutContext(), item)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestStorage_Get(t *testing.T) {
	t.Run("Data save in memory storage", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		st.CreateItem(context.Background(), domain.Item{Name: "Alex"})

		item, err := st.GetItem(context.Background(), 1)
		if item.ID != 1 || item.Name != "Alex" {
			t.Fatalf("expected item id= 1, name: Alex; got: id: %v, name: %v", item.ID, item.Name)
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("Returns error ErrNotFound", func(t *testing.T) {
		st := storage.NewMemoryStorage()

		_, err := st.GetItem(context.Background(), 1)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected error ErrNotFound, got: %v", err)
		}
	})

	t.Run("Context Canceled returns context.Canceled", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		st.CreateItem(context.Background(), domain.Item{Name: "Alex"})

		_, err := st.GetItem(CanceledContext(), 1)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Context timeout returns context.DeadlineExceeded", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		st.CreateItem(context.Background(), domain.Item{Name: "Alex"})
		_, err := st.GetItem(TimeoutContext(), 1)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestStorage_Delete(t *testing.T) {
	t.Run("Success delete", func(t *testing.T) {
		st := storage.NewMemoryStorage()
		st.CreateItem(context.Background(), domain.Item{Name: "Alex"})

		st.DeleteItem(context.Background(), 1)
		_, err := st.GetItem(context.Background(), 1)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got error: %v", err)
		}
	})

	t.Run("Returns error ErrNotFound", func(t *testing.T) {
		st := storage.NewMemoryStorage()

		err := st.DeleteItem(context.Background(), 1)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected error ErrNotFound, got: %v", err)
		}
	})

	t.Run("Context Canceled returns context.Canceled", func(t *testing.T) {
		st := storage.NewMemoryStorage()

		err := st.DeleteItem(CanceledContext(), 1)

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Context timeout returns context.DeadlineExceeded", func(t *testing.T) {
		st := storage.NewMemoryStorage()

		err := st.DeleteItem(TimeoutContext(), 1)

		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestStorage_ConcurrentCRUD(t *testing.T) {
	st := storage.NewMemoryStorage()
	const n = 1000

	// Этап 1: CREATE
	ids := make([]int, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			item, err := st.CreateItem(context.Background(), domain.Item{
				Name: fmt.Sprintf("item-%d", i),
			})
			if err != nil {
				t.Errorf("Create error: %v", err)
				return
			}
			ids[i] = item.ID
		}(i)
	}

	wg.Wait() // дождались всех созданных элементов

	// Этап 2: GET
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			item, err := st.GetItem(context.Background(), ids[i])
			if err != nil {
				t.Errorf("Get error for id=%d: %v", ids[i], err)
				return
			}
			if item.ID != ids[i] {
				t.Errorf("Get mismatch: expected id=%d, got id=%d", ids[i], item.ID)
			}
		}(i)
	}

	// Этап 3: DELETE
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			err := st.DeleteItem(context.Background(), ids[i])
			if err != nil {
				t.Errorf("Delete error for id=%d: %v", ids[i], err)
			}
		}(i)
	}

	wg.Wait() // ждём завершения GET и DELETE
}
