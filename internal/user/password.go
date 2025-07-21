package user

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
)

func GetPasswordHash(dtb *sql.DB, username string) (string, error) {
	var passwordHash string
	err := dtb.QueryRow("SELECT password_hash FROM users WHERE username = $1;", username).Scan(&passwordHash)
	if err != nil {
		return "", fmt.Errorf("error while selecting the password_hash: %v", err)
	}
	return passwordHash, nil
}

func HashPassword(password string) (string, error) {
	salt := "|3%$cris2QJlfs|R"

	hasher := sha256.New()
	_, err := hasher.Write([]byte(password + salt))
	if err != nil {
		return "", fmt.Errorf("error while hashing the password")
	}
	hashedBytes := hasher.Sum(nil)

	hashedPassword := hex.EncodeToString(hashedBytes)
	return hashedPassword, nil
}

func CheckPassword(password, passwordHash string) (bool, error) {
	hashed, err := HashPassword(password)
	if err != nil {
		return false, err
	}
	return hashed == passwordHash, nil
}
