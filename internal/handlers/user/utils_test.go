package user_test

import (
	"encoding/json"
	"io"
	"marketplace/internal/cards"
	"marketplace/internal/datastore"
	"marketplace/internal/handlers"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetUserHandler(t *testing.T) *uhd.UserHandler {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		t.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	usr := user.NewDBRepo(dtb)
	cards := cards.NewDBRepo(dtb)
	userHandler := &uhd.UserHandler{
		UserRepo:  usr,
		CardsRepo: cards,
	}
	return userHandler
}

func HandleMethodNotAllowed(t *testing.T, resp *http.Response) {
	code := resp.StatusCode
	if code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusMethodNotAllowed, code)
	}
}

func HandleBadReq(t *testing.T, resp *http.Response, expected string) {
	code := resp.StatusCode
	if code != http.StatusBadRequest {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusBadRequest, code)
	}

	HandleError(t, resp, expected)
}

func HandleError(t *testing.T, resp *http.Response, expected string) {
	if mime := resp.Header.Get("Content-Type"); mime != "application/json" {
		t.Errorf("Заголовок Content-Type должен иметь MIME-тип application/json, но имеет %s", mime)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка чтения тела ответа: %v", err)
		return
	}

	var errResp handlers.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatalf("Ошибка десериализации тела ответа сервера: %v", err)
	}

	result := errResp.Reason
	if result != expected {
		t.Errorf("Ожидалось %s, но получено %s", expected, result)
	}
}

func CheckCodeAndMime(t *testing.T, rr *httptest.ResponseRecorder) {
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, rr.Code)
	}

	if mime := rr.Header().Get("Content-Type"); mime != "application/json" {
		t.Errorf("Заголовок Content-Type должен иметь MIME-тип application/json, но имеет %s", mime)
	}
}
