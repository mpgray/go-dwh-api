package main

import (
	"fmt"
	"go-hoa-api/app"
	"go-hoa-api/controllers"
	u "go-hoa-api/utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/v1/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/v1/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/v1/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/v1/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8989" //localhost
	}

	u.Log.Info("Connected on port " + port)

	handler := corsConfig().Handler(router)
	err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:8989/api
	if err != nil {
		u.Log.Fatal(fmt.Sprint(err))
	}

	defer u.Log.Infof("**Golang Backend API for Driveway Home Started Successfully**")
}

func corsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://localhost",
			"https://127.0.0.1",
			"http://localhost",
			"http://127.0.0.1",
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
}
