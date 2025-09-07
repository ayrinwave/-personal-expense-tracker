package transport

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// userCtxKey - это ключ для хранения ID пользователя в контексте запроса.
// Используем собственный тип, чтобы избежать коллизий ключей.
type userCtxKey string

const UserIDKey userCtxKey = "userID"

// AuthMiddleware создает middleware для проверки JWT токена.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Получаем заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// 2. Проверяем, что заголовок имеет формат "Bearer <token>"
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := headerParts[1]

			// 3. Парсим и валидируем токен
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Убеждаемся, что используется правильный алгоритм подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// 4. Извлекаем ID пользователя из токена и добавляем в контекст
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userIDFloat, ok := claims["sub"].(float64); ok {
					userID := int(userIDFloat)
					// Создаем новый контекст с ID пользователя
					ctx := context.WithValue(r.Context(), UserIDKey, userID)
					// Вызываем следующий обработчик с новым контекстом
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		})
	}
}
