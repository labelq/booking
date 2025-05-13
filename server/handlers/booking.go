package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server/utils"
)

// BookParkingSpot - обработчик для бронирования парковочного места
func BookParkingSpot(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка авторизации (передаем userID через JWT)
		tokenString := r.Header.Get("Authorization")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Преобразуем "sub" (userID) из claims в тип int с безопасной проверкой
		userID, ok := claims["sub"].(float64) // jwt.MapClaims возвращает float64
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		// Преобразуем userID в int
		userIDInt := int(userID)

		// Получаем данные о бронировании
		var bookingData struct {
			ParkingSpot int `json:"parking_spot"`
		}
		err = json.NewDecoder(r.Body).Decode(&bookingData)
		if err != nil {
			http.Error(w, "Invalid booking data", http.StatusBadRequest)
			return
		}

		// Начинаем транзакцию
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Database transaction error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback() // Откатываем транзакцию, если что-то пошло не так

		// Проверяем, забронировано ли это место
		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM bookings WHERE parking_spot = $1", bookingData.ParkingSpot).Scan(&count)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Error(w, "Parking spot already booked", http.StatusBadRequest)
			return
		}

		// Бронирование места
		_, err = tx.Exec("INSERT INTO bookings (user_id, parking_spot) VALUES ($1, $2)", userIDInt, bookingData.ParkingSpot)
		if err != nil {
			http.Error(w, "Error while booking", http.StatusInternalServerError)
			return
		}

		// Подтверждаем транзакцию
		if err := tx.Commit(); err != nil {
			http.Error(w, "Error while committing transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Booking successful!")
	}
}
