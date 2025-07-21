package utils

import (
	"database/sql"
	"fmt"
)

func GetUsername(dtb *sql.DB, userID string) (*string, error) {
	var username string
	err := dtb.QueryRow("SELECT username FROM users WHERE id = $1;", userID).Scan(&username)
	if err != nil {
		return nil, fmt.Errorf("error while selecting the username of the user: %v", err)
	}
	return &username, nil
}
