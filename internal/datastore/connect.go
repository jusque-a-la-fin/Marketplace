package datastore

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func CreateNewDB() (*sql.DB, error) {
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	database := os.Getenv("DATABASE_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	dtb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database: %v", err)
	}

	dtb.SetMaxOpenConns(25)
	dtb.SetMaxIdleConns(25)
	dtb.SetConnMaxLifetime(5 * time.Minute)

	return dtb, nil
}
