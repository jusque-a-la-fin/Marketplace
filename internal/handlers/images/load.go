package images

import (
	"encoding/json"
	"log"
	hdr "marketplace/internal/handlers"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

const (
	UUIDRE   string = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	mimeType string = "image/jpeg"
)

type LoadImageRequest struct {
	Image  []byte `json:"image"`
	UserID string `json:"user_id"`
}

var regex = regexp.MustCompile(UUIDRE)

// LoadImage загружает изображение
func (hnd *ImagesHandler) LoadImage(wrt http.ResponseWriter, rqt *http.Request) {
	var lrq LoadImageRequest
	err := json.NewDecoder(rqt.Body).Decode(&lrq)
	if err != nil {
		errSend := hdr.SendBadReq(wrt, "wrong request body")
		log.Println("ERROR!!!", err.Error())
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	if !regex.MatchString(lrq.UserID) {
		errSend := hdr.SendBadReq(wrt, "Строка не соответствует формату UUID")
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	parsedUUID, err := uuid.Parse(lrq.UserID)
	if err != nil {
		errSend := hdr.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	imageName, err := hnd.ImagesRepo.LoadImage(mimeType, lrq.Image, parsedUUID)
	if err != nil {
		errSend := hdr.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)

	resp := struct {
		ImageName string `json:"image_name"`
	}{
		ImageName: *imageName,
	}

	errJSON := json.NewEncoder(wrt).Encode(resp)
	if errJSON != nil {
		log.Printf("error while sending response body: %v\n", errJSON)
	}
}
