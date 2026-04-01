package auth

import (
	"fmt"
)

// AuthService — интерфейс (контракт)
// нужен для абстракции и подмены реализации (например, в тестах)
type AuthService interface {
	Validate(token string) (string, error)
	Login(login, password string) (string, error)
}

// authService — реализация AuthService
// использует JWTValidation для проверки токена
type authService struct {
	jwt *JWTValidation
}

// конструктор сервиса авторизации
// возвращает интерфейс, а не конкретную реализацию (скрываем детали)
func NewAuthService(j *JWTValidation) AuthService {
	return &authService{jwt: j}
}

// Validate — делегирует проверку токена в JWT слой
func (a *authService) Validate(token string) (string, error) {
	return a.jwt.Parse(token)
}

// Login - проверяет логин, пароль и возвращает токен
func (a *authService) Login(login, password string) (string, error) {
	// проверяем логин и пароль
	if login != "admin" || password != "123" {
		return "", fmt.Errorf("invalid credentials")

	}

	// присваиваем id
	userID := "1"

	// возвращаем токен
	return a.jwt.Generate(userID)
}
