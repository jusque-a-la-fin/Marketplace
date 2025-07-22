package user

import (
	"encoding/json"
	"log"
	"marketplace/internal/cards"
	"marketplace/internal/handlers"
	"marketplace/internal/middleware"
	"marketplace/internal/token"
	"net/http"
	"strconv"
	"strings"
)

func (hnd *UserHandler) GetCards(wrt http.ResponseWriter, rqt *http.Request) {
	queryParams := rqt.URL.Query()
	page := 1
	if pageParam := queryParams.Get("page"); pageParam != "" {
		if pageInt, err := strconv.Atoi(pageParam); err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	perPage := 20
	if perPageParam := queryParams.Get("per_page"); perPageParam != "" {
		if perPageInt, err := strconv.Atoi(perPageParam); err == nil && perPageInt > 0 {
			perPage = perPageInt
		}
	}

	offset := (page - 1) * perPage
	sortBy := strings.ToLower(queryParams.Get("sort_by"))
	switch sortBy {
	case "price":
	default:
		sortBy = "created_at"
	}

	order := strings.ToLower(queryParams.Get("order"))
	if order != "asc" {
		order = "desc"
	}

	var priceMin, priceMax *float64

	if priceMinParam := queryParams.Get("price_min"); priceMinParam != "" {
		if priceMinFloat, err := strconv.ParseFloat(priceMinParam, 64); err == nil {
			priceMin = &priceMinFloat
		}
	}
	if priceMaxParam := queryParams.Get("price_max"); priceMaxParam != "" {
		if priceMaxFloat, err := strconv.ParseFloat(priceMaxParam, 64); err == nil {
			priceMax = &priceMaxFloat
		}
	}

	authVal := rqt.Context().Value(middleware.KeyIsAuthenticated)
	isAuthenticated, _ := authVal.(bool)

	var username *string = nil
	if isAuthenticated {
		usernameStr, err := token.GetPayload(rqt)
		if err != nil {
			errSend := handlers.SendUnauthorized(wrt, err.Error())
			if errSend != nil {
				log.Printf("error while sending the unauthorized error message: %v\n", errSend)
			}
			return
		}
		username = &usernameStr
	}

	params := &cards.QueryParams{
		PerPage:  perPage,
		Offset:   offset,
		SortBy:   sortBy,
		Order:    order,
		PriceMin: priceMin,
		PriceMax: priceMax,
		Username: username,
	}

	cards, err := hnd.CardsRepo.GetCards(params)
	if err != nil {
		errSend := handlers.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
	errJSON := json.NewEncoder(wrt).Encode(cards)
	if errJSON != nil {
		log.Printf("error while sending response body: %v\n", errJSON)
	}
}
