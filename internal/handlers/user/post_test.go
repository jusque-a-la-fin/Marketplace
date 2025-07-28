package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"marketplace/internal/datastore"
	ihd "marketplace/internal/handlers/images"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupTestServerForPostACard(t *testing.T) (*httptest.Server, *uhd.UserHandler) {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("error while connecting to the database: %v", err)
	}

	var uhr = GetUserHandler(t)
	var ihr = GetImagesHandler(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-in", uhr.SignIn).Methods("POST")
	rtr.HandleFunc("/post-a-card", middleware.RequireAuth(uhr.PostACard, dtb, true)).Methods("POST")
	rtr.HandleFunc("/images/create", ihr.CreateImage).Methods("GET")
	rtr.HandleFunc("/images", ihr.LoadImage).Methods("POST")
	path := fmt.Sprintf("/images/{name:image%v\\.jpeg}", ihd.UUIDRE)
	rtr.HandleFunc(path, ihr.GetImage).Methods("GET")

	ts := httptest.NewServer(rtr)
	t.Cleanup(ts.Close)
	return ts, uhr
}

// TestPostACard тестирует сценарий создания нового объявления
func TestPostACard(t *testing.T) {
	ts, uhr := setupTestServerForPostACard(t)
	name := "тест на сценарий создания объявления"
	t.Run(name, func(t *testing.T) {
		auth := uhd.AuthRequest{Username: "user1", Password: "W#_?e9o!m+B>tk7j"}
		data, err := json.Marshal(auth)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
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

		token := resp.Header.Get("Authorization")

		imageURL := getImageURL(t, ts, uhr, auth.Username)
		card := uhd.PostACardRequest{Title: "title1", Text: "text1", ImageURL: imageURL, Price: "1000"}
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
