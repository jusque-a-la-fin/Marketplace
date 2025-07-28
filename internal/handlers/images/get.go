package images

import (
	"fmt"
	"log"
	hdr "marketplace/internal/handlers"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

const ImagePattern string = `(image%s)\.jpeg`

var regexpr = regexp.MustCompile(fmt.Sprintf(ImagePattern, UUIDRE))

// GetImage получает изображение
func (hnd *ImagesHandler) GetImage(wrt http.ResponseWriter, rqt *http.Request) {
	vars := mux.Vars(rqt)
	rawImageName := vars["name"]
	matches := regexpr.FindStringSubmatch(rawImageName)

	var imageName *string
	if len(matches) > 1 {
		imageName = &matches[1]
	} else {
		errSend := hdr.SendBadReq(wrt, "неправильное имя изображения")
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	image, err := hnd.ImagesRepo.GetImage(*imageName)
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
