package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/utils"
	"strconv"
	"time"
)

type BookingRequest struct {
	ParkingSpot int    `json:"parkingSpot"`
	CarNumber   string `json:"carNumber"`
	Hours       int    `json:"hours"`
}

type BookingResponse struct {
	ID         int       `json:"id"`
	ReservedAt time.Time `json:"reservedAt"`
	EndTime    time.Time `json:"endTime"`
	Message    string    `json:"message"`
}

func BookParkingSpot(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Получаем токен
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		log.Printf("Token received: %v", tokenString != "")

		// Парсим токен
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Printf("Token parsing error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Получаем ID пользователя
		userID, ok := claims["sub"].(float64)
		if !ok {
			log.Printf("User ID conversion error. Claims: %v", claims)
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}
		userIDInt := int(userID)
		log.Printf("User ID: %d", userIDInt)

		// Декодируем тело запроса
		var bookingData BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&bookingData); err != nil {
			log.Printf("Request body decode error: %v", err)
			http.Error(w, "Invalid booking data", http.StatusBadRequest)
			return
		}
		log.Printf("Booking data received: %+v", bookingData)

		available, err := IsParkingSpotAvailable(db, bookingData.ParkingSpot)
		if err != nil {
			log.Printf("Error checking parking spot availability: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if !available {
			log.Printf("Parking spot %d is not available", bookingData.ParkingSpot)
			http.Error(w, "Parking spot is not available", http.StatusConflict)
			return
		}

		// Валидация данных
		if bookingData.ParkingSpot < 1 || bookingData.ParkingSpot > 16 {
			log.Printf("Invalid parking spot: %d", bookingData.ParkingSpot)
			http.Error(w, "Invalid parking spot number", http.StatusBadRequest)
			return
		}
		if bookingData.Hours < 1 {
			log.Printf("Invalid hours: %d", bookingData.Hours)
			http.Error(w, "Hours must be greater than 0", http.StatusBadRequest)
			return
		}
		if bookingData.CarNumber == "" {
			log.Printf("Empty car number")
			http.Error(w, "Car number is required", http.StatusBadRequest)
			return
		}

		// Начинаем транзакцию
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Transaction begin error: %v", err)
			http.Error(w, "Database transaction error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Проверяем, не занято ли место
		var count int
		err = tx.QueryRow(`
            SELECT COUNT(*)
            FROM bookings
            WHERE parking_spot = $1
            AND reserved_at + (hours * interval '1 hour') > NOW()
        `, bookingData.ParkingSpot).Scan(&count)

		if err != nil {
			log.Printf("Checking occupied spots error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			log.Printf("Parking spot %d is already occupied", bookingData.ParkingSpot)
			http.Error(w, "Parking spot is already booked", http.StatusConflict)
			return
		}

		// Создаем бронирование
		reservedAt := time.Now()
		var bookingID int
		err = tx.QueryRow(`
            INSERT INTO bookings (user_id, parking_spot, car_number, reserved_at, hours)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `, userIDInt, bookingData.ParkingSpot, bookingData.CarNumber, reservedAt, bookingData.Hours).Scan(&bookingID)

		if err != nil {
			log.Printf("Insert booking error: %v", err)
			http.Error(w, "Error while booking", http.StatusInternalServerError)
			return
		}

		// Подтверждаем транзакцию
		if err := tx.Commit(); err != nil {
			log.Printf("Transaction commit error: %v", err)
			http.Error(w, "Error while committing transaction", http.StatusInternalServerError)
			return
		}

		// Формируем ответ
		endTime := reservedAt.Add(time.Duration(bookingData.Hours) * time.Hour)
		response := BookingResponse{
			ID:         bookingID,
			ReservedAt: reservedAt,
			EndTime:    endTime,
			Message:    "Booking successful!",
		}

		log.Printf("Booking successful: %+v", response)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Response encoding error: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func IsParkingSpotAvailable(db *sql.DB, spotNumber int) (bool, error) {
	// Проверяем, не заблокировано ли место
	var isBlocked bool
	err := db.QueryRow(`
        SELECT is_blocked
        FROM blocked_spots
        WHERE spot_number = $1 AND is_blocked = true
    `, spotNumber).Scan(&isBlocked)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if isBlocked {
		return false, nil
	}

	// Проверяем, не занято ли место
	var count int
	err = db.QueryRow(`
        SELECT COUNT(*)
        FROM bookings
        WHERE parking_spot = $1
        AND reserved_at + (hours * interval '1 hour') > NOW()
    `, spotNumber).Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func GetAllBookings(db *sql.DB) http.HandlerFunc {
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

		// Получаем все активные бронирования
		rows, err := db.Query(`
            SELECT id, parking_spot, car_number, reserved_at, hours
            FROM bookings
            WHERE status = 'active'
            ORDER BY reserved_at DESC
        `)
		if err != nil {
			log.Printf("Database query error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookings []map[string]interface{}
		for rows.Next() {
			var id, parkingSpot, hours int
			var carNumber string
			var reservedAt time.Time

			if err := rows.Scan(&id, &parkingSpot, &carNumber, &reservedAt, &hours); err != nil {
				log.Printf("Row scan error: %v", err)
				continue
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

func CancelBooking(db *sql.DB) http.HandlerFunc {
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

		// Получаем ID бронирования из URL
		vars := mux.Vars(r)
		bookingID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid booking ID", http.StatusBadRequest)
			return
		}

		// Отменяем бронирование
		result, err := db.Exec(`
            UPDATE bookings
            SET status = 'cancelled'
            WHERE id = $1 AND status = 'active'
        `, bookingID)

		if err != nil {
			log.Printf("Database update error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Booking not found or already cancelled", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Booking cancelled successfully",
		})
	}
}

func ToggleSpotBlock(db *sql.DB) http.HandlerFunc {
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

		// Получаем номер места из тела запроса
		var req struct {
			SpotNumber int `json:"spotNumber"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверяем существование записи для данного места
		var exists bool
		var isCurrentlyBlocked bool
		err = db.QueryRow(`
            SELECT EXISTS(
                SELECT 1 FROM blocked_spots WHERE spot_number = $1
            ),
            COALESCE(
                (SELECT is_blocked FROM blocked_spots WHERE spot_number = $1),
                false
            )
        `, req.SpotNumber).Scan(&exists, &isCurrentlyBlocked)

		if err != nil {
			log.Printf("Database query error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Выполняем блокировку/разблокировку
		if exists {
			_, err = db.Exec(`
                UPDATE blocked_spots
                SET is_blocked = NOT is_blocked,
                    blocked_at = CURRENT_TIMESTAMP
                WHERE spot_number = $1
            `, req.SpotNumber)
		} else {
			_, err = db.Exec(`
                INSERT INTO blocked_spots (spot_number, is_blocked, blocked_at)
                VALUES ($1, true, CURRENT_TIMESTAMP)
            `, req.SpotNumber)
		}

		if err != nil {
			log.Printf("Database update error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Parking spot status updated successfully",
		})
	}
}
