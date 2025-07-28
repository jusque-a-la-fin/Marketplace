package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"marketplace/internal/cards"
	"marketplace/internal/handlers"
	img "marketplace/internal/images"
	"marketplace/internal/token"
	"marketplace/internal/utils"
	"net/http"
	"strconv"
)

const (
	// minTitleLen — минимальная длина заголовка
	minTitleLen int = 2
	// maxTitleLen — максимальная длина заголовка
	maxTitleLen int = 100
	// minTextLen — минимальная длина текста объявления
	minTextLen = 5
	// maxTextLen — максимальная длина текста объявления
	maxTextLen = 4000
	// maxPriceValue — максимальная цена
	maxPriceValue float64 = 1000000000000
)

// запрос с данными для создания объявления
type PostACardRequest struct {
	// Title — заголовок
	Title string `json:"title"`
	// Text — текст объявления
	Text string `json:"text"`
	// ImageURL — ссылка на изображение
	ImageURL string `json:"image_url"`
	// Price — цена
	Price string `json:"price"`
}

// PostACard создает новое объявление
func (hnd *UserHandler) PostACard(wrt http.ResponseWriter, rqt *http.Request) {
	var prq PostACardRequest
	err := json.NewDecoder(rqt.Body).Decode(&prq)
	if err != nil {
		errSend := handlers.SendBadReq(wrt, "wrong request body")
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	check := utils.CheckLen(prq.Title, "недостаточная", "превышена допустимая", "заголовка", "заголовок", minTitleLen, maxTitleLen)
	if check != "" {
		errSend := handlers.SendBadReq(wrt, check)
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	check = utils.CheckLen(prq.Title, "недостаточная", "превышена допустимая", "текста объявления", "текст объявления", minTextLen, maxTextLen)
	if check != "" {
		errSend := handlers.SendBadReq(wrt, check)
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	if err := validatePrice(prq.Price); err != nil {
		errSend := handlers.SendBadReq(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	if err := validateImage(prq.ImageURL); err != nil {
		errSend := handlers.SendBadReq(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the bad request message: %v\n", errSend)
		}
		return
	}

	priceFloat64, err := strconv.ParseFloat(prq.Price, 64)
	if err != nil {
		errStr := fmt.Errorf("error while conversion from string to float64: %v", err)
		errSend := handlers.SendInternalServerError(wrt, errStr.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	crd := &cards.CardInput{Title: prq.Title, Text: prq.Text, ImageURL: prq.ImageURL, Price: priceFloat64}
	username, err := token.GetPayload(rqt)
	if err != nil {
		errSend := handlers.SendUnauthorized(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the unauthorized error message: %v\n", errSend)
		}
		return
	}

	userID, err := hnd.UserRepo.GetUserID(username)
	if err != nil {
		errSend := handlers.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	card, err := hnd.CardsRepo.PostACard(crd, userID)
	if err != nil {
		errSend := handlers.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(card)
	if errJSON != nil {
		log.Printf("error while sending response body: %v\n", errJSON)
	}
}

// validatePrice валидирует цену
func validatePrice(priceStr string) error {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("error while conversion from string to float64: %v", err)
	}
	if price <= 0 {
		return fmt.Errorf("цена должна быть выше 0")
	}
	if price > maxPriceValue {
		return fmt.Errorf("цена не может быть выше %.0f", maxPriceValue)
	}
	return nil
}

// validateImage валидирует изображение
func validateImage(imageURL string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to issue a GET request: %v", err)
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка получения изображения, которое должно быть создано: %v", err)
	}

	size := len(imageData)
	if size == 0 {
		return fmt.Errorf("empty image data")
	}
	if size > img.MaxImageBytes {
		return fmt.Errorf("изображение превышает максимальный размер изображения = %d байтов, размер полученного изображения = %d байтов", img.MaxImageBytes, size)
	}

	// Определение content-type по первым 512 байтам
	firstBytesLen := 512
	if size < firstBytesLen {
		firstBytesLen = size
	}

	contentType := http.DetectContentType(imageData[:firstBytesLen])
	switch contentType {
	case "image/jpeg", "image/png":
	default:
		return fmt.Errorf("неизвестный content-type %q", contentType)
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("cannot decode image: %v", err)
	}

	width, height := cfg.Width, cfg.Height
	if width < img.MinImageDim || height < img.MinImageDim {
		return fmt.Errorf("недостаточное разрешение изображения. Разрешение полученного изображения = (%dx%d), минимальное возможное разрешение изображения = (%dx%d)", width, height, img.MinImageDim, img.MinImageDim)
	}
	if width > img.MaxImageDim || height > img.MaxImageDim {
		return fmt.Errorf("превышено максимально возможное разрешение изображения. Разрешение полученного изображения = (%dx%d), максимально возможное разрешение изображения = (%dx%d)", width, height, img.MaxImageDim, img.MaxImageDim)
	}

	ratio := float64(width) / float64(height)
	if ratio < img.MinAspectRatio || ratio > img.MaxAspectRatio {
		return fmt.Errorf("неправильное соотношение ширины к высоте полученного изображения = %.2f. Диапазон допустимого соотношения ширины к высоте изображения = %.2f–%.2f", ratio, img.MinAspectRatio, img.MaxAspectRatio)
	}
	return nil
}
