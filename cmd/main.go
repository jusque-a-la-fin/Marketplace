package main

import (
	"fmt"
	"log"
	"marketplace/internal/cards"
	"marketplace/internal/datastore"
	ihd "marketplace/internal/handlers/images"
	uhd "marketplace/internal/handlers/user"
	"marketplace/internal/images"
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

	images := images.NewDBRepo(dtb)
	imagesHandler := &ihd.ImagesHandler{
		ImagesRepo: images,
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/sign-in", userHandler.SignIn).Methods("POST")
	rtr.HandleFunc("/sign-up", userHandler.SignUp).Methods("POST")
	rtr.HandleFunc("/post-a-card", middleware.RequireAuth(userHandler.PostACard, dtb, true)).Methods("POST")
	rtr.HandleFunc("/get-cards", middleware.RequireAuth(userHandler.GetCards, dtb, false)).Methods("GET")
	rtr.HandleFunc("/images/create", imagesHandler.CreateImage).Methods("GET")
	rtr.HandleFunc("/images", imagesHandler.LoadImage).Methods("POST")
	path := fmt.Sprintf("/images/{name:image%s\\.jpeg}", ihd.UUIDRE)
	rtr.HandleFunc(path, imagesHandler.GetImage).Methods("GET")
	port := os.Getenv("SERVER_PORT")
	addr := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(addr, rtr); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
