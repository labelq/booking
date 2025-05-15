package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "server/utils"

)

// GetOccupiedSpots возвращает список занятых парковочных мест
func GetOccupiedSpots(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        // Проверяем токен
        tokenString := r.Header.Get("Authorization")
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        if _, err := utils.ParseToken(tokenString); err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Получаем занятые места
        rows, err := db.Query(`
            SELECT parking_spot
            FROM bookings
            WHERE reserved_at + (hours * interval '1 hour') > NOW()
        `)
        if err != nil {
            log.Printf("Database query error: %v", err)
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        occupiedSpots := []int{}
        for rows.Next() {
            var spot int
            if err := rows.Scan(&spot); err != nil {
                log.Printf("Error scanning row: %v", err)
                http.Error(w, "Error reading data", http.StatusInternalServerError)
                return
            }
            occupiedSpots = append(occupiedSpots, spot)
        }

        if err = rows.Err(); err != nil {
            log.Printf("Error iterating rows: %v", err)
            http.Error(w, "Error reading data", http.StatusInternalServerError)
            return
        }

        // Отправляем ответ
        response := map[string]interface{}{
            "occupiedSpots": occupiedSpots,
        }

        if err := json.NewEncoder(w).Encode(response); err != nil {
            log.Printf("Error encoding response: %v", err)
            http.Error(w, "Error encoding response", http.StatusInternalServerError)
            return
        }
    }
}