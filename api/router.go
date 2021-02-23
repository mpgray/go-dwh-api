package api

import (
	"go-dwh-api/controllers"
	"go-dwh-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

//api endpoints
const (
	// Do not need authorization
	newAccount = "/account/new"
	login      = "/account/login"
	refresh    = "/account/refresh"
	// Needs authorization
	logout = "/account/logout"

	newContact  = "/contact/new"
	getContact  = "/me/contact"
	getContacts = "/me/contacts"
)

// Router creates and serves the server
func Router() *gin.Engine {
	router := gin.Default()

	// Group that needs no authenticating, i.e. unauthenticated
	unauthenticated := router.Group("/api/v1")
	{
		unauthenticated.POST(login, controllers.Login)
		unauthenticated.POST(newAccount, controllers.CreateAccount)
		unauthenticated.POST(refresh, controllers.Refresh)
	}

	authenticated := router.Group("/api/v1/auth")
	// Group that requires an authenticated user
	authenticated.Use(models.TokenAuthenticator())
	{
		authenticated.POST(logout, controllers.Logout)
		authenticated.POST(newContact, controllers.CreateContact)
	}

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
