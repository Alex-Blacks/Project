package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTValidation — обёртка над JWT библиотекой
// отвечает только за парсинг и валидацию токена
type JWTValidation struct {
	secret []byte // секретный ключ для проверки подписи токена
}

// конструктор JWTValidation
func NewJWTValidation(secret []byte) *JWTValidation {
	return &JWTValidation{secret: secret}
}

// Parse — проверяет токен и извлекает userID
func (j *JWTValidation) Parse(tokenStr string) (string, error) {

	// парсим токен и одновременно проверяем подпись
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		// проверяем алгоритм подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		// возвращаем секретный ключ для проверки подписи
		return j.secret, nil
	})
	if err != nil {
		// ошибка парсинга (невалидная структура, подпись и т.д.)
		return "", err
	}

	// дополнительная проверка валидности токена
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// извлекаем claims (payload токена)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// если тип не тот — потенциальная проблема токена
		return "", fmt.Errorf("invalid claims")
	}

	// достаём user_id из payload
	userID, ok := claims["user_id"].(string)
	if !ok {
		// если нет user_id или он не string — токен некорректный
		return "", fmt.Errorf("user_id missing")
	}

	// возвращаем userID, который будет использоваться дальше в системе
	return userID, nil
}

func (j *JWTValidation) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}
