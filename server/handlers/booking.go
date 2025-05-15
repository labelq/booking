package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "server/utils"
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