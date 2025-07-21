package user

import (
	"database/sql"
	"fmt"
)

type UserDBRepository struct {
	dtb *sql.DB
}

func NewDBRepo(sdb *sql.DB) *UserDBRepository {
	return &UserDBRepository{dtb: sdb}
}

// GetUserID получает идентификатор пользователя
func (repo *UserDBRepository) GetUserID(username string) (string, error) {
	var userID string
	err := repo.dtb.QueryRow("SELECT id FROM users WHERE username = $1;", username).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("error while selecting the user id: %v", err)
	}
	return userID, nil
}
