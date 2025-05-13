package models

import (
	"database/sql"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	AccountType  string `json:"account_type"` // 'user' или 'admin'
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Создание нового пользователя
func CreateUser(db *sql.DB, user *User) error {
	query := "INSERT INTO users (email, password_hash, account_type) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, user.Email, user.PasswordHash, user.AccountType)
	return err
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
