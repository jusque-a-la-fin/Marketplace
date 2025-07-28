package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"marketplace/internal/cards"
	"marketplace/internal/datastore"
	ihd "marketplace/internal/handlers/images"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupTestServerForGetCards(t *testing.T) (*httptest.Server, *uhd.UserHandler) {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		t.Fatalf("error while connecting to the database: %v", err)
	}

	userHandler := GetUserHandler(t)
	var ihr = GetImagesHandler(t)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-in", userHandler.SignIn).Methods("POST")
	rtr.HandleFunc("/sign-up", userHandler.SignUp).Methods("POST")
	rtr.HandleFunc("/post-a-card", middleware.RequireAuth(userHandler.PostACard, dtb, true)).Methods("POST")
	rtr.HandleFunc("/get-cards", middleware.RequireAuth(userHandler.GetCards, dtb, false)).Methods("GET")
	rtr.HandleFunc("/images/create", ihr.CreateImage).Methods("GET")
	rtr.HandleFunc("/images", ihr.LoadImage).Methods("POST")
	path := fmt.Sprintf("/images/{name:image%v\\.jpeg}", ihd.UUIDRE)
	rtr.HandleFunc(path, ihr.GetImage).Methods("GET")

	ts := httptest.NewServer(rtr)
	t.Cleanup(ts.Close)
	return ts, userHandler
}

var testsForAuthorized = map[string]struct {
	input  string
	result []cards.CardOutput
}{
	"url №1: /get-cards": {input: "/get-cards", result: []cards.CardOutput{
		{Title: "title10", Text: "text10", ImageURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3", IsOwned: true},
		{Title: "title9", Text: "text9", ImageURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3", IsOwned: true},
		{Title: "title8", Text: "text8", ImageURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1", IsOwned: false},
		{Title: "title7", Text: "text7", ImageURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1", IsOwned: false},
		{Title: "title6", Text: "text6", ImageURL: "https://www.example.com/images/image6.jpg", Price: 6000, Username: "user1", IsOwned: false},
		{Title: "title5", Text: "text5", ImageURL: "https://www.example.com/images/image5.jpg", Price: 5000, Username: "user1", IsOwned: false},
	}},
	"url №2: /get-cards?price_min=20&price_max=90&sort_by=price&order=asc": {input: "/get-cards?price_min=7000&price_max=10000&sort_by=price&order=asc",
		result: []cards.CardOutput{
			{Title: "title7", Text: "text7", ImageURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1", IsOwned: false},
			{Title: "title8", Text: "text8", ImageURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1", IsOwned: false},
			{Title: "title9", Text: "text9", ImageURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3", IsOwned: true},
			{Title: "title10", Text: "text10", ImageURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3", IsOwned: true},
		}},
}

var testsForUnauthorized = map[string]struct {
	input  string
	result []cards.CardOutput
}{
	"url №1: /get-cards": {input: "/get-cards", result: []cards.CardOutput{
		{Title: "title10", Text: "text10", ImageURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3"},
		{Title: "title9", Text: "text9", ImageURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3"},
		{Title: "title8", Text: "text8", ImageURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1"},
		{Title: "title7", Text: "text7", ImageURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1"},
		{Title: "title6", Text: "text6", ImageURL: "https://www.example.com/images/image6.jpg", Price: 6000, Username: "user1"},
		{Title: "title5", Text: "text5", ImageURL: "https://www.example.com/images/image5.jpg", Price: 5000, Username: "user1"},
	}},
	"url №2: /get-cards?price_min=20&price_max=90&sort_by=price&order=asc": {input: "/get-cards?price_min=7000&price_max=10000&sort_by=price&order=asc",
		result: []cards.CardOutput{
			{Title: "title7", Text: "text7", ImageURL: "https://www.example.com/images/image7.jpg", Price: 7000, Username: "user1"},
			{Title: "title8", Text: "text8", ImageURL: "https://www.example.com/images/image8.jpg", Price: 8000, Username: "user1"},
			{Title: "title9", Text: "text9", ImageURL: "https://www.example.com/images/image9.jpg", Price: 9000, Username: "user3"},
			{Title: "title10", Text: "text10", ImageURL: "https://www.example.com/images/image10.jpg", Price: 10000, Username: "user3"},
		}},
}

// TestGetCards тестирует сценарий получения ленты объявлений для зарегистрированного пользователя и для незарегистрированного пользователя
func TestGetCards(t *testing.T) {
	ts, uhr := setupTestServerForGetCards(t)
	auth := uhd.AuthRequest{
		Username: "user1",
		Password: "W#_?e9o!m+B>tk7j",
	}

	token := Authorize(t, ts, auth, "/sign-in")

	imageURL5 := getImageURL(t, ts, uhr, auth.Username)
	imageURL6 := getImageURL(t, ts, uhr, auth.Username)
	imageURL7 := getImageURL(t, ts, uhr, auth.Username)
	imageURL8 := getImageURL(t, ts, uhr, auth.Username)

	imageUrlsForUser1 := map[string]string{
		"imageURL5": imageURL5,
		"imageURL6": imageURL6,
		"imageURL7": imageURL7,
		"imageURL8": imageURL8,
	}

	PrepareTestsForUser1(testsForAuthorized, imageUrlsForUser1)
	PrepareTestsForUser1(testsForUnauthorized, imageUrlsForUser1)

	cardsToPost1 := []uhd.PostACardRequest{
		{Title: "title5", Text: "text5", ImageURL: imageURL5, Price: "5000"},
		{Title: "title6", Text: "text6", ImageURL: imageURL6, Price: "6000"},
		{Title: "title7", Text: "text7", ImageURL: imageURL7, Price: "7000"},
		{Title: "title8", Text: "text8", ImageURL: imageURL8, Price: "8000"},
	}

	PostCards(t, ts, cardsToPost1, token)

	auth = uhd.AuthRequest{
		Username: "user3",
		Password: "Q#_~s1o!m+B&t/9j0g{",
	}

	token = Authorize(t, ts, auth, "/sign-up")

	imageURL9 := getImageURL(t, ts, uhr, auth.Username)
	imageURL10 := getImageURL(t, ts, uhr, auth.Username)

	imageUrlsForUser3 := map[string]string{
		"imageURL9":  imageURL9,
		"imageURL10": imageURL10,
	}

	PrepareTestsForUser3(testsForAuthorized, imageUrlsForUser3)
	PrepareTestsForUser3(testsForUnauthorized, imageUrlsForUser3)

	cardsToPost2 := []uhd.PostACardRequest{
		{Title: "title9", Text: "text9", ImageURL: imageURL9, Price: "9000"},
		{Title: "title10", Text: "text10", ImageURL: imageURL10, Price: "10000"},
	}

	PostCards(t, ts, cardsToPost2, token)
	RunTests(t, ts, testsForAuthorized, token, true)
	RunTests(t, ts, testsForUnauthorized, token, false)
}

func PrepareTestsForUser1(tests map[string]struct {
	input  string
	result []cards.CardOutput
}, imageURLs map[string]string) {
	for key := range tests {
		for idx, cardOutput := range tests[key].result {
			switch cardOutput.Title {
			case "title5":
				tests[key].result[idx].ImageURL = imageURLs["imageURL5"]
			case "title6":
				tests[key].result[idx].ImageURL = imageURLs["imageURL6"]
			case "title7":
				tests[key].result[idx].ImageURL = imageURLs["imageURL7"]
			case "title8":
				tests[key].result[idx].ImageURL = imageURLs["imageURL8"]
			}
		}
	}
}

func PrepareTestsForUser3(tests map[string]struct {
	input  string
	result []cards.CardOutput
}, imageURLs map[string]string) {
	for key := range tests {
		for idx, cardOutput := range tests[key].result {
			switch cardOutput.Title {
			case "title9":
				tests[key].result[idx].ImageURL = imageURLs["imageURL9"]
			case "title10":
				tests[key].result[idx].ImageURL = imageURLs["imageURL10"]
			}
		}
	}
}

func PostCards(t *testing.T, ts *httptest.Server, cardsToPost []uhd.PostACardRequest, token string) {
	client := &http.Client{}
	for _, card := range cardsToPost {
		data, err := json.Marshal(card)
		if err != nil {
			t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
		}

		fullURL := fmt.Sprintf("%s%s", ts.URL, "/post-a-card")
		req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewReader(data))
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
}

func RunTests(t *testing.T, ts *httptest.Server, tests map[string]struct {
	input  string
	result []cards.CardOutput
}, token string, isAuthorized bool) {
	client := &http.Client{}
	for _, test := range tests {
		fullURL := fmt.Sprintf("%s%s", ts.URL, test.input)
		req, err := http.NewRequest(http.MethodGet, fullURL, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		if isAuthorized {
			req.Header.Set("Authorization", token)
		}
		resp, err := client.Do(req)
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

			if cards[i].ImageURL != expected.ImageURL {
				t.Errorf("Ожидался ImageURL: %q, но получен ImageURL: %q", expected.ImageURL, cards[i].ImageURL)
			}

			if cards[i].Price != expected.Price {
				t.Errorf("Ожидался Price: %f, но получен Price: %f", expected.Price, cards[i].Price)
			}

			if cards[i].Username != expected.Username {
				t.Errorf("Ожидался Username: %q, но получен Username: %q", expected.Username, cards[i].Username)
			}

			if isAuthorized {
				if cards[i].IsOwned != expected.IsOwned {
					t.Errorf("Ожидался IsOwned: %t, но получен IsOwned: %t", expected.IsOwned, cards[i].IsOwned)
				}
			}
		}
	}
}

func Authorize(t *testing.T, ts *httptest.Server, auth uhd.AuthRequest, path string) string {
	data, err := json.Marshal(auth)
	if err != nil {
		t.Fatalf("Ошибка сериализации тела запроса клиента: %v", err)
	}

	fullURL := fmt.Sprintf("%s%s", ts.URL, path)
	resp, err := http.Post(fullURL, "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to issue a POST request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
	}

	token := resp.Header.Get("Authorization")
	resp.Body.Close()
	return token
}
