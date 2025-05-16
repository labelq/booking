package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/utils"
	"strconv"
)

// Структуры для ответов
type AdminResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Обработчик для получения всех пользователей
func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем права администратора
		claims, err := utils.GetAndValidateTokenClaims(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accountType, ok := claims["account_type"].(string)
		if !ok || accountType != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Получаем список пользователей из базы данных
		rows, err := db.Query("SELECT id, email, account_type FROM users")
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var id int
			var email, accountType string
			if err := rows.Scan(&id, &email, &accountType); err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			users = append(users, map[string]interface{}{
				"id":           id,
				"email":        email,
				"account_type": accountType,
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": users,
		})
	}
}

// Обработчик для изменения роли пользователя
func UpdateUserRoleHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем права администратора
		claims, err := utils.GetAndValidateTokenClaims(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accountType, ok := claims["account_type"].(string)
		if !ok || accountType != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Получаем ID пользователя из URL
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Получаем новую роль из тела запроса
		var requestBody struct {
			Role string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Обновляем роль пользователя в базе данных
		_, err = db.Exec("UPDATE users SET account_type = $1 WHERE id = $2", requestBody.Role, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AdminResponse{
			Success: true,
			Message: "User role updated successfully",
		})
	}
}

// Обработчик для получения всех бронирований
func GetBookingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем права администратора
		claims, err := utils.GetAndValidateTokenClaims(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accountType, ok := claims["account_type"].(string)
		if !ok || accountType != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Получаем список активных бронирований
		rows, err := db.Query(`
            SELECT id, parking_spot, car_number, reserved_at, hours
            FROM bookings
            WHERE status = 'active'
        `)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookings []map[string]interface{}
		for rows.Next() {
			var id, parkingSpot, hours int
			var carNumber string
			var reservedAt string
			if err := rows.Scan(&id, &parkingSpot, &carNumber, &reservedAt, &hours); err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			bookings = append(bookings, map[string]interface{}{
				"id":           id,
				"parking_spot": parkingSpot,
				"car_number":   carNumber,
				"reserved_at":  reservedAt,
				"hours":        hours,
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"bookings": bookings,
		})
	}
}

// Обработчик для отмены бронирования
func CancelBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем права администратора
		claims, err := utils.GetAndValidateTokenClaims(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accountType, ok := claims["account_type"].(string)
		if !ok || accountType != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Получаем ID бронирования из URL
		vars := mux.Vars(r)
		bookingID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid booking ID", http.StatusBadRequest)
			return
		}

		// Отменяем бронирование
		_, err = db.Exec("UPDATE bookings SET status = 'cancelled' WHERE id = $1", bookingID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AdminResponse{
			Success: true,
			Message: "Booking cancelled successfully",
		})
	}
}

// Обработчик для получения заблокированных мест
func GetBlockedSpots(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем права администратора
		claims, err := utils.GetAndValidateTokenClaims(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !isAdminFromClaims(claims) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT spot_number
			FROM blocked_spots
			WHERE is_blocked = true
		`)

		if err != nil {
			log.Printf("Database query error: %v", err) // Добавьте это
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		blockedSpots := []int{} // Инициализируем пустым массивом
		for rows.Next() {
			var spotNumber int
			if err := rows.Scan(&spotNumber); err != nil {
				log.Printf("Row scan error: %v", err)
				continue
			}
			blockedSpots = append(blockedSpots, spotNumber)
		}

		// Всегда возвращаем объект с массивом, даже если он пустой
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"blockedSpots": blockedSpots,
		})
	}
}

// Обработчик для блокировки/разблокировки места
func ToggleSpotBlockHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем права администратора
		claims, err := utils.GetAndValidateTokenClaims(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accountType, ok := claims["account_type"].(string)
		if !ok || accountType != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Получаем номер места из тела запроса
		var requestBody struct {
			SpotNumber int `json:"spotNumber"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверяем существование записи для данного места
		var exists bool
		var isCurrentlyBlocked bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM blocked_spots WHERE spot_number = $1), COALESCE((SELECT is_blocked FROM blocked_spots WHERE spot_number = $1), false)",
			requestBody.SpotNumber).Scan(&exists, &isCurrentlyBlocked)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if exists {
			// Обновляем существующую запись
			_, err = db.Exec("UPDATE blocked_spots SET is_blocked = NOT is_blocked WHERE spot_number = $1",
				requestBody.SpotNumber)
		} else {
			// Создаем новую запись
			_, err = db.Exec("INSERT INTO blocked_spots (spot_number, is_blocked) VALUES ($1, true)",
				requestBody.SpotNumber)
		}

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AdminResponse{
			Success: true,
			Message: "Spot status updated successfully",
		})
	}
}
