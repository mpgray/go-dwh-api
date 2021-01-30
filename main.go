package main

import (
	"fmt"
	"go-hoa-api/app"
	"go-hoa-api/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
