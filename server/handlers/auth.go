package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/models"
	"server/utils"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        var userData models.User

        if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
            log.Printf("Error decoding request: %v", err)
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        // Получаем пароль до того, как он будет захэширован
        password := userData.PasswordHash // временно храним пароль

        // Хэширование пароля
        hashedPassword, err := utils.HashPassword(password) // используем исходный пароль
        if err != nil {
            log.Printf("Error hashing password: %v", err)
            http.Error(w, "Could not hash password", http.StatusInternalServerError)
            return
        }

        userData.PasswordHash = hashedPassword
        userData.AccountType = "user"

		// Создаем пользователя и получаем его ID
		userID, err := models.CreateUser(db, &userData) // Измените CreateUser, чтобы возвращать ID
		if err != nil {
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// Генерируем токен с правильным ID
		token, err := utils.GenerateToken(userID)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":    userID,
				"email": userData.Email,
			},
		})
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        var loginData models.LoginData

        if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
            log.Printf("Login: error decoding request: %v", err)
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        log.Printf("Login attempt for email: %s", loginData.Email)

        user, err := models.FindUserByEmail(db, loginData.Email)
        if err != nil {
            log.Printf("Login: error finding user: %v", err)
            if err == sql.ErrNoRows {
                http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            } else {
                http.Error(w, "Database error", http.StatusInternalServerError)
            }
            return
        }

        if !utils.CheckPassword(loginData.Password, user.PasswordHash) {
            log.Printf("Login: invalid password for user: %s", loginData.Email)
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        token, err := utils.GenerateToken(user.ID)
        if err != nil {
            log.Printf("Login: error generating token: %v", err)
            http.Error(w, "Could not create token", http.StatusInternalServerError)
            return
        }

        log.Printf("Login successful for user: %s", loginData.Email)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "token": token,
            "user": map[string]interface{}{
                "id":    user.ID,
                "email": user.Email,
            },
        })
    }
}