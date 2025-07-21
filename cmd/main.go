package main

import (
	"log"
	"marketplace/internal/cards"
	"marketplace/internal/datastore"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/middleware"
	"marketplace/internal/user"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	dtb, err := datastore.CreateNewDB()
	if err != nil {
		log.Fatalf("error while connecting to the database: %v", err)
	}

	usr := user.NewDBRepo(dtb)
	cards := cards.NewDBRepo(dtb)
	userHandler := &uhd.UserHandler{
		UserRepo:  usr,
		CardsRepo: cards,
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-in", userHandler.SignIn).Methods("POST")
	rtr.HandleFunc("/sign-up", userHandler.SignUp).Methods("POST")
	rtr.HandleFunc("/post-a-card", middleware.RequireAuth(userHandler.PostACard, dtb, true)).Methods("POST")
	rtr.HandleFunc("/get-cards", middleware.RequireAuth(userHandler.GetCards, dtb, false)).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	if err := http.ListenAndServe(":"+port, rtr); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
