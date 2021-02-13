package app

import (
	"fmt"
	"go-dwh-api/controllers"
	u "go-dwh-api/utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	NEW_ACCOUNT  = "/account/new"
	LOGIN        = "/account/login"
	NEW_CONTACT  = "/contact/new"
	GET_CONTACTS = "/me/contacts"
)

var router *mux.Router

func Serve() {
	router = mux.NewRouter()

	apiPath := os.Getenv("api_path")

	router.HandleFunc(apiPath+NEW_ACCOUNT, controllers.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc(apiPath+LOGIN, controllers.Authenticate).Methods(http.MethodPost)
	router.HandleFunc(apiPath+NEW_CONTACT, controllers.CreateContact).Methods(http.MethodPost)
	router.HandleFunc(apiPath+GET_CONTACTS, controllers.GetContactsFor).Methods(http.MethodGet) //  user/2/contacts

	router.Use(JwtAuthentication) //attach JWT auth middleware

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
