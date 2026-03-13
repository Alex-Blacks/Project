package domain

import (
	"context"
)

// Интерфейс для storage
type Storage interface {
	CreateItem(ctx context.Context, item Item) (Item, error) // Создать элемент
	GetItem(ctx context.Context, id int) (Item, error)       // Отправить элемент
	DeleteItem(ctx context.Context, id int) error            // Удалить элемент
}

// Интерфейс для gRPC
type TaskService interface {
	Create(ctx context.Context, name string) (Item, error) // Создать элемент
	Get(ctx context.Context, id int) (Item, error)         // Отправить элемент
	Delete(ctx context.Context, id int) error              // Удалить элемент
}
