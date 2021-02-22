package api

import (
	"go-dwh-api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

//api endpoints
const (
	newAccount  = "/api/v1/account/new"
	login       = "/api/v1/account/login"
	logout      = "/api/v1/account/logout"
	refresh     = "/api/v1/account/refresh"
	newContact  = "/api/v1/contact/new"
	getContact  = "/api/v1/me/contact"
	getContacts = "/api/v1/me/contacts"
)

// Router creates and serves the server
func Router() *gin.Engine {
	router := gin.Default()

	router.POST(login, controllers.Login)
	router.POST(newAccount, controllers.CreateAccount)
	router.POST(refresh, controllers.Refresh)
	router.POST(logout, controllers.Logout)
	router.POST(newContact, controllers.CreateContact)

	//	router.HandleFunc(apiPath+newAccount, controllers.CreateAccount).Methods(http.MethodPost)
	//	router.HandleFunc(apiPath+login, controllers.Authenticate).Methods(http.MethodPost)
	//  router.HandleFunc(apiPath+newContact, controllers.CreateContact).Methods(http.MethodPost)
	//  router.HandleFunc(apiPath+getContact, controllers.GetContact).Methods(http.MethodPost)
	//  router.HandleFunc(apiPath+getContacts, controllers.GetContactsFor).Methods(http.MethodGet) //  user/2/contacts

	//	router.Use(JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	return router
}

// CorsConfig controls the cross origin
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
