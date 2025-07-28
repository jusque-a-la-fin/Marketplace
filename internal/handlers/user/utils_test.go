package user_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"marketplace/internal/cards"
	"marketplace/internal/datastore"
	"marketplace/internal/handlers"
	ihd "marketplace/internal/handlers/images"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/images"
	"marketplace/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ConnectToDB(t *testing.T) *sql.DB {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		t.Fatalf("error while connecting to database: %v", err)
	}
	return dtb
}

func GetUserHandler(t *testing.T) *uhd.UserHandler {
	dtb := ConnectToDB(t)
	usr := user.NewDBRepo(dtb)
	cards := cards.NewDBRepo(dtb)
	userHandler := &uhd.UserHandler{
		UserRepo:  usr,
		CardsRepo: cards,
	}
	return userHandler
}

func GetImagesHandler(t *testing.T) *ihd.ImagesHandler {
	dtb := ConnectToDB(t)
	images := images.NewDBRepo(dtb)
	imagesHandler := &ihd.ImagesHandler{
		ImagesRepo: images,
	}
	return imagesHandler
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
		t.Fatalf("error while reading response body: %v", err)
		return
	}

	var errResp handlers.ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		t.Fatalf("error while deserialization response body from server: %v", err)
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

func getImageURL(t *testing.T, ts *httptest.Server, uhr *uhd.UserHandler, username string) string {
	resp, err := http.Get(ts.URL + "/images/create")
	if err != nil {
		t.Fatalf("failed to issue a GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
	}

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Ошибка получения изображения, которое должно быть создано: %v", err)
	}

	userID, err := uhr.UserRepo.GetUserID(username)
	if err != nil {
		t.Fatalf("Ошибка получения ID пользователя: %v", err)
	}

	loadImageRequest := struct {
		Image  []byte `json:"image"`
		UserID string `json:"user_id"`
	}{
		Image:  image,
		UserID: userID,
	}

	data, err := json.Marshal(loadImageRequest)
	if err != nil {
		t.Fatalf("error while serialization response body for client: %v", err)
	}

	resp, err = http.Post(ts.URL+"/images", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("failed to issue a POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался код состояния ответа: %d, но получен: %d", http.StatusOK, resp.StatusCode)
	}

	loadResp := struct {
		ImageName string `json:"image_name"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&loadResp)
	if err != nil {
		t.Fatalf("error while deserialization response body from server: %v", err)
	}

	imageURL := fmt.Sprintf("%s/images/%s.jpeg", ts.URL, loadResp.ImageName)
	return imageURL
}
