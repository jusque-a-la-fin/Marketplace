package user

import (
	"marketplace/internal/cards"
	"marketplace/internal/user"
)

type UserHandler struct {
	UserRepo  user.UserRepo
	CardsRepo cards.CardsRepo
}
