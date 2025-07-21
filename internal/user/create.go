package user

import (
	"database/sql"
	"fmt"
)

// SignUp записывает нового пользователя в базу данных
func CreateUser(dtb *sql.DB, user *User) (*User, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	var userID string
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id;`
	err = dtb.QueryRow(query, user.Username, hashedPassword).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("error while inserting a new user: %v", err)
	}

	thisUser := User{Username: user.Username, Password: hashedPassword}
	return &thisUser, nil
}
