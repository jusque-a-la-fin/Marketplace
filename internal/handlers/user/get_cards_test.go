package user_test

import (
	"bytes"
	"encoding/json"
	"marketplace/internal/cards"
	"marketplace/internal/datastore"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupTestServer(t *testing.T) (*httptest.Server, *uhd.UserHandler) {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		t.Fatalf("error while connecting to the database: %v", err)
	}

	userHandler := GetUserHandler(t)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-in", userHandler.SignIn).Methods("POST")
	rtr.HandleFunc("/sign-up", userHandler.SignUp).Methods("POST")
	rtr.HandleFunc("/post-a-card", middleware.RequireAuth(userHandler.PostACard, dtb, true)).Methods("POST")
	rtr.HandleFunc("/get-cards", middleware.RequireAuth(userHandler.GetCards, dtb, false)).Methods("GET")

	ts := httptest.NewServer(rtr)
	t.Cleanup(ts.Close)
	return ts, userHandler
}

var testsForAthourized = map[string]struct {
	input  string
	result []cards.CardOutput
}{
	"url №1: /get-cards": {input: "/get-cards", result: []cards.CardOutput{
		{Title: "title10", Text: "text10", PictureURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3", IsOwned: true},
		{Title: "title9", Text: "text9", PictureURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3", IsOwned: true},
		{Title: "title8", Text: "text8", PictureURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1", IsOwned: false},
		{Title: "title7", Text: "text7", PictureURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1", IsOwned: false},
		{Title: "title6", Text: "text6", PictureURL: "https://www.example.com/images/image6.jpg", Price: 6000, Username: "user1", IsOwned: false},
		{Title: "title5", Text: "text5", PictureURL: "https://www.example.com/images/image5.jpg", Price: 5000, Username: "user1", IsOwned: false},
	}},
	"url №2: /get-cards?price_min=20&price_max=90&sort_by=price&order=asc": {input: "/get-cards?price_min=7000&price_max=10000&sort_by=price&order=asc",
		result: []cards.CardOutput{
			{Title: "title7", Text: "text7", PictureURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1", IsOwned: false},
			{Title: "title8", Text: "text8", PictureURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1", IsOwned: false},
			{Title: "title9", Text: "text9", PictureURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3", IsOwned: true},
			{Title: "title10", Text: "text10", PictureURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3", IsOwned: true},
		}},
}

var testsForUnathourized = map[string]struct {
	input  string
	result []cards.CardOutput
}{
	"url №1: /get-cards": {input: "/get-cards", result: []cards.CardOutput{
		{Title: "title10", Text: "text10", PictureURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3", IsOwned: false},
		{Title: "title9", Text: "text9", PictureURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3", IsOwned: false},
		{Title: "title8", Text: "text8", PictureURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1", IsOwned: false},
		{Title: "title7", Text: "text7", PictureURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1", IsOwned: false},
		{Title: "title6", Text: "text6", PictureURL: "https://www.example.com/images/image6.jpg", Price: 6000, Username: "user1", IsOwned: false},
		{Title: "title5", Text: "text5", PictureURL: "https://www.example.com/images/image5.jpg", Price: 5000, Username: "user1", IsOwned: false},
	}},
	"url №2: /get-cards?price_min=20&price_max=90&sort_by=price&order=asc": {input: "/get-cards?price_min=7000&price_max=10000&sort_by=price&order=asc",
		result: []cards.CardOutput{
			{Title: "title7", Text: "text7", PictureURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1", IsOwned: false},
			{Title: "title8", Text: "text8", PictureURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1", IsOwned: false},
			{Title: "title9", Text: "text9", PictureURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3", IsOwned: false},
			{Title: "title10", Text: "text10", PictureURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3", IsOwned: false},
		}},
}

func TestGetCards(t *testing.T) {
	ts, _ := setupTestServer(t)
	auth := uhd.AuthRequest{
		Username: "user1",
		Password: "W#_?e9o!m+B>tk7j",
	}

	data, err := json.Marshal(auth)
	if err != nil {
		t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
	}

	resp, err := http.Post(ts.URL+"/sign-in", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to issue a POST: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
	}

	token := resp.Header.Get("Authorization")
	resp.Body.Close()

	client := &http.Client{}

	cardsToPost1 := []uhd.PostACardRequest{
		{Title: "title5", Text: "text5", PictureURL: "https://www.example.com/images/image5.jpg", Price: "5000"},
		{Title: "title6", Text: "text6", PictureURL: "https://www.example.com/images/image6.jpg", Price: "6000"},
		{Title: "title7", Text: "text7", PictureURL: "https://www.example.com/images/image7.jpg", Price: "7000"},
		{Title: "title8", Text: "text8", PictureURL: "https://www.example.com/images/image8.jpg", Price: "8000"},
	}

	for _, card := range cardsToPost1 {
		data, err = json.Marshal(card)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, ts.URL+"/post-a-card", bytes.NewReader(data))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make a request: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 on post-a-card, cards %d", resp.StatusCode)
		}
		defer resp.Body.Close()
	}

	auth = uhd.AuthRequest{
		Username: "user3",
		Password: "Q#_~s1o!m+B&t/9j0g{",
	}
	data, err = json.Marshal(auth)
	if err != nil {
		t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
	}

	resp, err = http.Post(ts.URL+"/sign-up", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to issue a POST: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
	}
	token = resp.Header.Get("Authorization")
	resp.Body.Close()

	cardsToPost2 := []uhd.PostACardRequest{
		{Title: "title9", Text: "text9", PictureURL: "https://www.example.com/images/image9.jpg", Price: "9000"},
		{Title: "title10", Text: "text10", PictureURL: "https://www.example.com/images/image10.jpg", Price: "10000"},
	}

	for _, card := range cardsToPost2 {
		data, err = json.Marshal(card)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, ts.URL+"/post-a-card", bytes.NewReader(data))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make a request: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
		}
		defer resp.Body.Close()
	}

	for _, test := range testsForAthourized {
		req, err := http.NewRequest(http.MethodGet, ts.URL+test.input, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", token)
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make a request: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
		}

		var cards []cards.CardOutput
		if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
			t.Fatalf("Ошибка десериализации ответа сервера: %v", err)
		}
		resp.Body.Close()

		if len(cards) < len(test.result) {
			t.Fatalf("expected at least %d cards, cards %d", len(test.result), len(cards))
		}

		for i := range test.result {
			expected := test.result[i]
			if cards[i].Title != expected.Title {
				t.Errorf("Ожидался Title: %q, но получен Title: %q", expected.Title, cards[i].Title)
			}

			if cards[i].Text != expected.Text {
				t.Errorf("Ожидался Text: %q, но получен Text: %q", expected.Text, cards[i].Text)
			}

			if cards[i].PictureURL != expected.PictureURL {
				t.Errorf("Ожидался PictureURL: %q, но получен PictureURL: %q", expected.PictureURL, cards[i].PictureURL)
			}

			if cards[i].Price != expected.Price {
				t.Errorf("Ожидался Price: %f, но получен Price: %f", expected.Price, cards[i].Price)
			}

			if cards[i].Username != expected.Username {
				t.Errorf("Ожидался Username: %q, но получен Username: %q", expected.Username, cards[i].Username)
			}

			if cards[i].IsOwned != expected.IsOwned {
				t.Errorf("Ожидался IsOwned: %t, но получен IsOwned: %t", expected.IsOwned, cards[i].IsOwned)
			}
		}
	}

	for _, test := range testsForUnathourized {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+test.input, nil)
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("get-cards request failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 on get-cards, cards %d", resp.StatusCode)
		}

		var cards []cards.CardOutput
		if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
			t.Fatalf("decoding get-cards response: %v", err)
		}
		resp.Body.Close()

		if len(cards) < len(test.result) {
			t.Fatalf("expected at least %d cards, cards %d", len(test.result), len(cards))
		}

		for i := range test.result {
			expected := test.result[i]
			if cards[i].Title != expected.Title {
				t.Errorf("Ожидался Title: %q, но получен Title: %q", expected.Title, cards[i].Title)
			}

			if cards[i].Text != expected.Text {
				t.Errorf("Ожидался Text: %q, но получен Text: %q", expected.Text, cards[i].Text)
			}

			if cards[i].PictureURL != expected.PictureURL {
				t.Errorf("Ожидался PictureURL: %q, но получен PictureURL: %q", expected.PictureURL, cards[i].PictureURL)
			}

			if cards[i].Price != expected.Price {
				t.Errorf("Ожидался Price: %f, но получен Price: %f", expected.Price, cards[i].Price)
			}

			if cards[i].Username != expected.Username {
				t.Errorf("Ожидался Username: %q, но получен Username: %q", expected.Username, cards[i].Username)
			}

			if cards[i].IsOwned != expected.IsOwned {
				t.Errorf("Ожидался IsOwned: %t, но получен IsOwned: %t", expected.IsOwned, cards[i].IsOwned)
			}
		}
	}

}
