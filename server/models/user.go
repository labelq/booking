package models

import (
	"database/sql"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	AccountType  string `json:"account_type"` // 'user' или 'admin'
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(db *sql.DB, user *User) (int, error) {
    var id int
    err := db.QueryRow(`
        INSERT INTO users (email, password_hash, account_type)
        VALUES ($1, $2, $3)
        RETURNING id
    `, user.Email, user.PasswordHash, user.AccountType).Scan(&id)

    if err != nil {
        return 0, err
    }

    return id, nil
}

// Проверка существующего пользователя по email
func FindUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User
	query := "SELECT id, email, password_hash, account_type FROM users WHERE email = $1"
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.AccountType)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
