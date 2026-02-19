package domain

import (
	"context"
)

type Storage interface {
	CreateItem(ctx context.Context, item Item) (Item, error) // Создать элемент
	GetItem(ctx context.Context, id int) (Item, error)       // Отправить элемент
	DeleteItem(ctx context.Context, id int) error            // Удалить элемент
}
