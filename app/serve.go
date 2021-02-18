package app

import (
	"go-dwh-api/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Full path is /api/v1/{const}
const (
	newAccount  = "/account/new"
	login       = "/account/login"
	newContact  = "/contact/new"
	getContact  = "/me/contact"
	getContacts = "/me/contacts"
)

var router *mux.Router

// Serve creates and serves the server
func Serve() *mux.Router {
	router = mux.NewRouter()

	apiPath := os.Getenv("api_path")

	router.HandleFunc(apiPath+newAccount, controllers.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc(apiPath+login, controllers.Authenticate).Methods(http.MethodPost)
	router.HandleFunc(apiPath+newContact, controllers.CreateContact).Methods(http.MethodPost)
	router.HandleFunc(apiPath+getContact, controllers.GetContact).Methods(http.MethodPost)
	router.HandleFunc(apiPath+getContacts, controllers.GetContactsFor).Methods(http.MethodGet) //  user/2/contacts

	router.Use(JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	return router
}

func CorsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://localhost",
			"https://127.0.0.1",
			"http://localhost",
			"http://127.0.0.1",
			"https://localhost:4200",
			"http://localhost:4200",
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
