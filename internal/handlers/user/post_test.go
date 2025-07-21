package user_test

import (
	"bytes"
	"encoding/json"
	uhd "marketplace/internal/handlers/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Интеграционный тест на сценарий создания объявления
func TestPostACard(t *testing.T) {
	var uhr = GetUserHandler(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-in", uhr.SignIn).Methods("POST")
	rtr.HandleFunc("/post-a-card", uhr.PostACard).Methods("POST")

	ts := httptest.NewServer(rtr)
	t.Cleanup(ts.Close)

	name := "тест на сценарий создания объявления"
	t.Run(name, func(t *testing.T) {
		auth := uhd.AuthRequest{Username: "user1", Password: "W#_?e9o!m+B>tk7j"}
		data, err := json.Marshal(auth)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		resp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("Failed to issue a POST: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
		}

		token := resp.Header.Get("Authorization")

		card := uhd.PostACardRequest{Title: "title1", Text: "text1", PictureURL: "https://www.example.com/images/image1.jpg", Price: "1000"}
		data, err = json.Marshal(card)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, ts.URL+"/post-a-card", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make a request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
		}
	})
}
