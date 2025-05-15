package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "server/utils"
    "strconv"
    "strings"
)

type UserResponse struct {
    ID          int    `json:"id"`      // Исправлен тег json
    Email       string `json:"email"`    // Исправлен тег json
    AccountType string `json:"account_type"` // Исправлен тег json
}

type RoleUpdateRequest struct {
    AccountType string `json:"account_type"` // Исправлен тег json
}

// Функция для проверки прав администратора
func isAdminFromClaims(claims map[string]interface{}) bool {
    accountType, ok := claims["account_type"].(string)
    if !ok || accountType != "admin" {
        return false
    }
    return true
}

func GetAllUsers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        claims, err := utils.GetAndValidateTokenClaims(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if !isAdminFromClaims(claims) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        // Получаем список пользователей
        rows, err := db.Query(`
            SELECT id, email, account_type
            FROM users
            ORDER BY id
        `)
        if err != nil {
            log.Printf("Database query error: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var users []UserResponse
        for rows.Next() {
            var user UserResponse
            if err := rows.Scan(&user.ID, &user.Email, &user.AccountType); err != nil {
                log.Printf("Row scan error: %v", err)
                continue
            }
            users = append(users, user)
        }

        json.NewEncoder(w).Encode(map[string]interface{}{
            "users": users,
        })
    }
}

func UpdateUserRole(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        claims, err := utils.GetAndValidateTokenClaims(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if !isAdminFromClaims(claims) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        parts := strings.Split(r.URL.Path, "/")
        if len(parts) < 5 {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        userID, err := strconv.Atoi(parts[len(parts)-2])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        var req RoleUpdateRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        if req.AccountType != "admin" && req.AccountType != "user" {
            http.Error(w, "Invalid account type", http.StatusBadRequest)
            return
        }

        adminID := int(claims["sub"].(float64))
        if userID == adminID {
            http.Error(w, "Cannot change your own account type", http.StatusForbidden)
            return
        }

        result, err := db.Exec(`
            UPDATE users
            SET account_type = $1
            WHERE id = $2
        `, req.AccountType, userID)

        if err != nil {
            log.Printf("Database update error: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{
            "message": "User account type updated successfully",
        })
    }
}