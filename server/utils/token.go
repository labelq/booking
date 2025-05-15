package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
	"errors"
    "net/http"
    "strings"
)

var secretKey = []byte("your-secret-key") // Этот ключ можно хранить в env-файле

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(inputPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	return err == nil
}

func GenerateToken(userID int, accountType string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": userID,
        "account_type": accountType,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })

    return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not extract claims")
	}

	return claims, nil
}

func GetAndValidateTokenClaims(r *http.Request) (map[string]interface{}, error) {
    // Получаем токен из заголовка
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
        return nil, errors.New("no token provided")
    }

    // Убираем префикс "Bearer "
    if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
        tokenString = tokenString[7:]
    }

    // Парсим и проверяем токен
    claims, err := ParseToken(tokenString)
    if err != nil {
        return nil, err
    }

    return claims, nil
}