package utils

import (
	"database/sql"
	"fmt"
	"unicode/utf8"
)

func CheckUser(dtb *sql.DB, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);`
	err := dtb.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error while checking if the user exists: %v", err)
	}
	return exists, nil
}

func CheckLen(str, errStr1, errStr2, field1, field2 string, min, max int) string {
	length := utf8.RuneCountInString(str)
	if length < min || length > max {
		partErrStr := "ошибка:"
		errStr := fmt.Sprintf("длина %s -> %s должен содержать от %d до %d символов", field1, field2, min, max)

		if length < min {
			partErrStr = fmt.Sprintf("%s %s", partErrStr, errStr1)
			return fmt.Sprintf("%s %s", partErrStr, errStr)
		}

		if length > max {
			partErrStr = fmt.Sprintf("%s %s", partErrStr, errStr2)
			return fmt.Sprintf("%s %s", partErrStr, errStr)
		}
	}

	return ""
}
