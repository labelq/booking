package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

var secretKey = []byte("your-secret-key")

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		// Проверка на отсутствие токена
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			log.Println("Authorization token is missing") // Логируем отсутствие токена
			return
		}

		// Убираем префикс "Bearer " из токена
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Парсим и проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверка метода подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		// Обработка ошибок при парсинге токена
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Printf("Token parsing error: %v", err) // Логируем ошибку парсинга
			return
		}

		// Проверка, что токен действителен
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Println("Invalid token") // Логируем, что токен невалиден
			return
		}

		// Передаем управление следующему обработчику, если токен валиден
		next.ServeHTTP(w, r)
	})
}
