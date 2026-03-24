package domain

import "errors"

var (
	ErrEmptyName    = errors.New("empty name")            // пустое имя
	ErrInvalidValue = errors.New("invalid value")         // Ошибка значения
	ErrNotFound     = errors.New("not found")             // Нет данных
	ErrInternal     = errors.New("server internal error") // Ошибка сервера
	ErrBadRequest   = errors.New("bad request")           // ошибка запроса
)
