package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	hnd "marketplace/internal/handlers"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var urls = []string{"/sign-up", "/sign-in"}

func setupTestServerForAuthorization(t *testing.T) *httptest.Server {
	uhr := GetUserHandler(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-up", uhr.SignUp).Methods("POST")
	rtr.HandleFunc("/sign-in", uhr.SignIn).Methods("POST")

	ts := httptest.NewServer(rtr)
	t.Cleanup(ts.Close)
	return ts
}

// TestMethodNotAllowed тестирует некорректный метод запроса
func TestMethodNotAllowed(t *testing.T) {
	ts := setupTestServerForAuthorization(t)
	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			t.Parallel()

			fullURL := fmt.Sprintf("%s%s", ts.URL, url)
			resp, err := http.Get(fullURL)
			if err != nil {
				t.Fatalf("failed to issue a GET request: %v", err)
			}
			defer resp.Body.Close()

			HandleMethodNotAllowed(t, resp)
		})
	}
}

var testsBadRequest = map[string]struct {
	input  uhd.AuthRequest
	result hnd.ErrorResponse
}{
	"нет логина и пароля":                                  {input: uhd.AuthRequest{}, result: hnd.ErrorResponse{Reason: "ошибка: пользователь не отправил логин и пароль"}},
	"нет логина":                                           {input: uhd.AuthRequest{Password: "password"}, result: hnd.ErrorResponse{Reason: "ошибка: пользователь не отправил логин"}},
	"нет пароля":                                           {input: uhd.AuthRequest{Username: "user2"}, result: hnd.ErrorResponse{Reason: "ошибка: пользователь не отправил пароль"}},
	"недостаточная длина логина":                           {input: uhd.AuthRequest{Username: "aa", Password: "W#_?e9o!m+B>tk7j"}, result: hnd.ErrorResponse{Reason: "ошибка: недостаточная длина логина -> логин должен содержать от 3 до 20 символов"}},
	"превышена допустимая длина логина":                    {input: uhd.AuthRequest{Username: "aaaaaaaaaaaaaaaaaaaaa", Password: "W#_?e9o!m+B>tk7j"}, result: hnd.ErrorResponse{Reason: "ошибка: превышена допустимая длина логина -> логин должен содержать от 3 до 20 символов"}},
	"недопустимый символ в логине":                         {input: uhd.AuthRequest{Username: "aaaaaaaa#", Password: "W#_?e9o!m+B>tk7j"}, result: hnd.ErrorResponse{Reason: "ошибка: недопустимый/ые символ/ы в логине -> логин может содержать только буквы, цифры, '-' и '_'"}},
	"такой логин уже занят":                                {input: uhd.AuthRequest{Username: "user1", Password: "W#_?e9o!m+B>tk7j"}, result: hnd.ErrorResponse{Reason: "ошибка: такой логин уже занят"}},
	"недостаточная длина пароля":                           {input: uhd.AuthRequest{Username: "user2", Password: "W#_?"}, result: hnd.ErrorResponse{Reason: "ошибка: недостаточная длина пароля -> пароль должен содержать от 8 до 30 символов"}},
	"превышена допустимая длина пароля":                    {input: uhd.AuthRequest{Username: "user2", Password: "W#_?e9o!m+B>tk7jw#_?e9o!m+b>tk7j"}, result: hnd.ErrorResponse{Reason: "ошибка: превышена допустимая длина пароля -> пароль должен содержать от 8 до 30 символов"}},
	"в пароле отсутствует хотя бя одна заглавная буква":    {input: uhd.AuthRequest{Username: "user2", Password: "w#_?e9o!m+b>tk7j"}, result: hnd.ErrorResponse{Reason: "ошибка: в пароле отсутствует хотя бя одна заглавная буква"}},
	"в пароле отсутствует хотя бя одна строчная буква":     {input: uhd.AuthRequest{Username: "user2", Password: "W#_?E9O!M+B>TK7J"}, result: hnd.ErrorResponse{Reason: "ошибка: в пароле отсутствует хотя бя одна строчная буква"}},
	"в пароле отсутствует хотя бя одна цифра":              {input: uhd.AuthRequest{Username: "user2", Password: "W#_?eo!m+B>tkj"}, result: hnd.ErrorResponse{Reason: "ошибка: в пароле отсутствует хотя бя одна цифра"}},
	"в пароле отсутствует хотя бя один специальный символ": {input: uhd.AuthRequest{Username: "user2", Password: "We9omBtk7j"}, result: hnd.ErrorResponse{Reason: "ошибка: в пароле отсутствует хотя бя один специальный символ"}},
}

func setupTestServerForSignUp(t *testing.T) *httptest.Server {
	uhr := GetUserHandler(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-up", uhr.SignUp).Methods("POST")

	ts := httptest.NewServer(rtr)
	t.Cleanup(ts.Close)
	return ts
}

// TestBadRequest тестирует некорректные запросы на регистрацию
func TestBadRequest(t *testing.T) {
	ts := setupTestServerForSignUp(t)

	// некорректные параметры тела запроса
	for name, testBR := range testsBadRequest {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(testBR.input)
			if err != nil {
				t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
			}

			fullURL := fmt.Sprintf("%s%s", ts.URL, "/sign-up")
			resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(data))
			if err != nil {
				t.Fatalf("failed to issue a POST request: %v", err)
			}
			defer resp.Body.Close()

			HandleBadReq(t, resp, testBR.result.Reason)
		})
	}
}

// TestSignUp тестирует сценарий регистрации нового пользователя
func TestSignUp(t *testing.T) {
	ts := setupTestServerForSignUp(t)
	name := "тест на сценарий регистрации"
	t.Run(name, func(t *testing.T) {
		auth := uhd.AuthRequest{Username: "user2", Password: "&N^_?e9G!m+B>[k3a"}
		data, err := json.Marshal(auth)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		fullURL := fmt.Sprintf("%s%s", ts.URL, "/sign-up")
		resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("failed to issue a POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
		}

		if mime := resp.Header.Get("Content-Type"); mime != "application/json" {
			t.Errorf("Заголовок Content-Type должен иметь MIME-тип application/json, но имеет %s", mime)
		}

		var usr user.User
		if err := json.NewDecoder(resp.Body).Decode(&usr); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		hashedPassword, err := user.HashPassword(auth.Password)
		if err != nil {
			t.Fatalf("error while hashing the password: %v", err)
		}

		if usr.Username != auth.Username && usr.Password != hashedPassword {
			t.Error("Логин, полученный от сервера, не совпадает с логином, отправленным пользователем. \n Также хэшированный пароль, полученный от сервера, не совпадает с отправленным пользователем паролем, который был захэширован.")
		}

		if usr.Username != auth.Username {
			t.Error("Логин, полученный от сервера, не совпадает с логином, отправленным пользователем.")
		}

		if usr.Password != hashedPassword {
			t.Error("Хэшированный пароль, полученный от сервера, не совпадает с отправленным пользователем паролем, который был захэширован.")
		}
	})
}

func setupTestServerForSignIn(t *testing.T) *httptest.Server {
	uhr := GetUserHandler(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-in", uhr.SignIn).Methods("POST")

	ts := httptest.NewServer(rtr)
	t.Cleanup(ts.Close)
	return ts
}

// TestSignIn тестирует сценарий авторизации зарегистрированного пользователя
func TestSignIn(t *testing.T) {
	ts := setupTestServerForSignIn(t)
	name := "тест на сценарий авторизации зарегистрированного пользователя"
	t.Run(name, func(t *testing.T) {
		auth := uhd.AuthRequest{Username: "user1", Password: "W#_?e9o!m+B>tk7j"}
		data, err := json.Marshal(auth)
		if err != nil {
			t.Fatalf("error while serialization: %v", err)
		}

		fullURL := fmt.Sprintf("%s%s", ts.URL, "/sign-in")
		resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("failed to issue a POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
		}
	})
}
