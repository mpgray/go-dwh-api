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

	u.Log.Info(port)

	handler := cors.Default().Handler(router)     // TODO: configure cors to allow only acceptable domains
	err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:8989/api
	if err != nil {
		u.Log.Error(fmt.Sprint(err))
	}
}
