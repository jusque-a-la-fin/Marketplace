package token

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/utils"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ExampleTokenSecret = []byte("ExampleTokenSecret")
	ErrNoToken         = errors.New("no token was in the request")
)

func CreateJWTtoken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": username,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + 1300,
	})
	tokenString, err := token.SignedString(ExampleTokenSecret)
	return tokenString, err
}

func Check(rqt *http.Request, dtb *sql.DB) (bool, error) {
	username, err := GetPayload(rqt)
	if err != nil {
		return false, err
	}

	exists, err := utils.CheckUser(dtb, username)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetPayload(rqt *http.Request) (string, error) {
	inToken := rqt.Header.Get("Authorization")
	if inToken == "" {
		return "", ErrNoToken
	}

	hashSecretGetter := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("expected another signing method")
		}
		return ExampleTokenSecret, nil
	}

	token, errJwt := jwt.Parse(inToken, hashSecretGetter)
	if errJwt != nil {
		return "", errJwt
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error while fetching the payload")
	}
	claims := payload["user"].(map[string]any)
	return claims["username"].(string), nil
}
