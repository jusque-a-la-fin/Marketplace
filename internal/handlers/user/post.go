package user

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"marketplace/internal/cards"
	"marketplace/internal/handlers"
	"marketplace/internal/token"
	"marketplace/internal/utils"
	"net/http"
	"net/url"
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
	// maxImageBytes — максимальный размер файла изображения в байтах (2 Мбайта)
	maxImageBytes = 2000000
	// minImageDim — минимальное количество пикселей по каждой стороне изображения
	minImageDim = 500
	// maxImageDim — максимальное количество пикселей по каждой стороне изображения
	maxImageDim = 2000
	// minAspectRatio — минимально допустимое соотношение ширины к высоте изображения
	minAspectRatio = 0.8
	// maxAspectRatio — максимально допустимое соотношение ширины к высоте изображения
	maxAspectRatio = 1.2
)

// запрос с данными для создания объявления
type PostACardRequest struct {
	// Title — заголовок
	Title string `json:"title"`
	// Text — текст объявления
	Text string `json:"text"`
	// PictureURL — адрес изображения
	PictureURL string `json:"picture_url"`
	// Price — цена
	Price string `json:"price"`
}

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
			log.Printf("ошибка отправки Bad Request сообщения %v\n", errSend)
		}
		return
	}

	check = utils.CheckLen(prq.Title, "недостаточная", "превышена допустимая", "текста объявления", "текст объявления", minTextLen, maxTextLen)
	if check != "" {
		errSend := handlers.SendBadReq(wrt, check)
		if errSend != nil {
			log.Printf("ошибка отправки Bad Request сообщения %v\n", errSend)
		}
		return
	}

	if err := validatePrice(prq.Price); err != nil {
		errSend := handlers.SendBadReq(wrt, err.Error())
		if errSend != nil {
			log.Printf("ошибка отправки Bad Request сообщения %v\n", errSend)
		}
		return
	}

	// if err := validateImage(prq.PictureURL); err != nil {
	// 	errSend := handlers.SendBadReq(wrt, err.Error())
	// 	if errSend != nil {
	// 		log.Printf("ошибка отправки Bad Request сообщения %v\n", errSend)
	// 	}
	// 	return
	// }

	priceFloat64, err := strconv.ParseFloat(prq.Price, 64)
	if err != nil {
		errStr := fmt.Errorf("ошибка при преобразовании из string в float64: %v", err)
		errSend := handlers.SendInternalServerError(wrt, errStr.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	crd := &cards.CardInput{Title: prq.Title, Text: prq.Text, PictureURL: prq.PictureURL, Price: priceFloat64}
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

// validatePrice проверяет цену
func validatePrice(priceStr string) error {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("price must be a valid number")
	}
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if price > maxPriceValue {
		return fmt.Errorf("price must be at most %.0f", maxPriceValue)
	}
	return nil
}

// validateImage проверяет изображение на соответствие формату
func validateImage(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL")
	}

	resp, err := http.Head(u.String())
	if err != nil {
		return fmt.Errorf("cannot HEAD URL")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("image URL returned status %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "image/jpeg" && ct != "image/png" {
		return fmt.Errorf("unsupported content-type %s", ct)
	}

	if cl := resp.Header.Get("Content-Length"); cl != "" {
		size, err := strconv.Atoi(cl)
		if err != nil {
			return fmt.Errorf("invalid content-length")
		}
		if size > maxImageBytes {
			return fmt.Errorf("file too large: %d bytes (max %d)", size, maxImageBytes)
		}
	}

	resp2, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("cannot download image")
	}
	defer resp2.Body.Close()

	limited := io.LimitReader(resp2.Body, maxImageBytes+1)
	cfg, _, err := image.DecodeConfig(limited)
	if err != nil {
		return fmt.Errorf("cannot decode image: %v", err)
	}

	w, h := cfg.Width, cfg.Height
	if w < minImageDim || h < minImageDim {
		return fmt.Errorf("image dimensions too small: %dx%d (min %dx%d)", w, h, minImageDim, minImageDim)
	}
	if w > maxImageDim || h > maxImageDim {
		return fmt.Errorf("image dimensions too large: %dx%d (max %dx%d)", w, h, maxImageDim, maxImageDim)
	}

	ratio := float64(w) / float64(h)
	if ratio < minAspectRatio || ratio > maxAspectRatio {
		return fmt.Errorf("invalid aspect ratio %.2f (allowed %.2f–%.2f)", ratio, minAspectRatio, maxAspectRatio)
	}

	return nil
}
