package api

import (
	"go-dwh-api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

//api endpoints
const (
	// Does not need authorization
	newUser = "/user/new"
	login   = "/account/login"
	// Needs authorization
	refresh     = "/account/refresh"
	logout      = "/account/logout"
	newContact  = "/contact/new"
	getContact  = "/my/contact"
	getContacts = "/my/contacts"
	getName     = "/my/contact/name"
	getNames    = "/my/contact/names"
)

// Router creates and serves the server
func Router() *gin.Engine {
	router := gin.Default()

	// Group that needs no authenticating, i.e. unauthenticated
	unauthenticated := router.Group("/api/v1")
	{
		unauthenticated.POST(login, controllers.Login)
		unauthenticated.POST(newUser, controllers.CreateUser)
	}

	authenticated := router.Group("/api/v1/auth")
	// Group that requires an authenticated user
	authenticated.Use(controllers.TokenAuthenticator())
	{
		authenticated.POST(logout, controllers.Logout)
		authenticated.POST(newContact, controllers.CreateContact)
		authenticated.POST(getContact, controllers.GetContact)
		authenticated.GET(getContacts, controllers.GetContactsFor)
		authenticated.POST(refresh, controllers.Refresh)
		authenticated.POST(getName, controllers.GetName)
		authenticated.GET(getNames, controllers.GetNamesFor)
	}

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
