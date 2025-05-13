package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/models"
	"server/utils"
)

// Регистрация пользователя
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	log.Println("dfhgsfnghn")
	return func(w http.ResponseWriter, r *http.Request) {
		var userData models.User

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		log.Println(r.Body)
		// Здесь passwordData - это поле, которое приходит от клиента, не PasswordHash
		password := userData.PasswordHash

		// Проверка, существует ли уже пользователь с таким email
		existingUser, err := models.FindUserByEmail(db, userData.Email)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if existingUser != nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}

		// Хэшируем пароль перед сохранением
		hashedPassword, err := utils.HashPassword(password) // Используем Password для хэширования
		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}

		// Сохраняем хэшированный пароль
		userData.PasswordHash = hashedPassword
		userData.AccountType = "user" // По умолчанию 'user', можно менять при регистрации администратора

		// Сохраняем нового пользователя в базе данных
		err = models.CreateUser(db, &userData)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// Генерируем JWT токен для нового пользователя
		token, err := utils.GenerateToken(userData.ID)
		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		// Отправляем токен в ответе
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginData models.LoginData
		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Ищем пользователя по email
		user, err := models.FindUserByEmail(db, loginData.Email)
		if err != nil {
			// Пользователь не найден или ошибка в запросе
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Проверка пароля
		if !utils.CheckPassword(loginData.Password, user.PasswordHash) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Генерация JWT токена
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		// Отправляем токен в ответе
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}
