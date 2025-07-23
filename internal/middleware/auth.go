package middleware

import (
	"context"
	"database/sql"
	"log"
	"marketplace/internal/token"
	"net/http"
)

type contextKey string

const (
	KeyIsAuthenticated = contextKey("isAuthenticated")
)

func RequireAuth(next http.HandlerFunc, dtb *sql.DB, isRequired bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		check, err := token.Check(r, dtb)
		if err != nil && (isRequired || err != token.ErrNoToken) {
			log.Printf("internal error during token check: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if isRequired && !check {
			log.Println("the token check has failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), KeyIsAuthenticated, check)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
