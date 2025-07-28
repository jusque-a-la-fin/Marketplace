package images

import (
	"log"
	hdr "marketplace/internal/handlers"
	"net/http"
	"strconv"
)

// CreateImage создает изображение
func (hnd *ImagesHandler) CreateImage(wrt http.ResponseWriter, rqt *http.Request) {
	image, err := hnd.ImagesRepo.CreateImage()
	if err != nil {
		errSend := hdr.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	wrt.Header().Set("Content-Type", "image/jpeg")
	wrt.Header().Set("Content-Length", strconv.Itoa(len(image)))
	wrt.WriteHeader(http.StatusOK)

	if _, err := wrt.Write(image); err != nil {
		log.Printf("error while writing image data: %v", err)
	}
}
