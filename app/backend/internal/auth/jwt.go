package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Подписать токен
func Sign(userID int64, role, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,                     // ваш ID
		"role": role,                       // роль
		"exp":  time.Now().Add(ttl).Unix(), // срок жизни (в секундах)
		"iat":  time.Now().Unix(),          // когда выдан
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

// Распарсить и проверить токен
// Parse разбирает JWT-токен, проверяет подпись и возвращает его содержимое (claims).
func Parse(tokenStr, secret string) (jwt.MapClaims, error) {
	// 1. Пробуем распарсить токен с нашими claims (payload).
	tok, err := jwt.ParseWithClaims(
		tokenStr,
		jwt.MapClaims{}, // говорим: хотим получить содержимое как map[string]interface{}
		func(t *jwt.Token) (interface{}, error) {
			// сюда библиотека придёт, чтобы спросить: "каким ключом проверить подпись?"
			return []byte(secret), nil // возвращаем наш секретный ключ
		},
	)

	// 2. Если произошла ошибка или токен невалидный — сразу выходим.
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid token")
	}

	// 3. Приводим tok.Claims к типу jwt.MapClaims (по сути это map[string]any).
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// 4. Возвращаем содержимое токена (например: sub, role, exp, iat).
	return claims, nil
}
