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
	refresh        = "/account/refresh"
	logout         = "/account/logout"
	newContact     = "/contact/new"
	newStatement   = "/statement/new"
	getStatement   = "/statement"
	getContact     = "/my/contact"
	getContacts    = "/my/contacts"
	searchContacts = "/my/contacts/search"
	getName        = "/my/contact/name"
	getNames       = "/my/contact/names"
	getAddress     = "/my/contact/address"
	getAddresses   = "/my/contact/addresses"
	getPhone       = "/my/contact/phone"
	getPhones      = "/my/contact/phones"
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
		authenticated.POST(refresh, controllers.Refresh)
		authenticated.POST(logout, controllers.Logout)
		authenticated.POST(getStatement, controllers.GetStatement)
		authenticated.POST(newStatement, controllers.CreateStatement)
		// All contact endpoints
		authenticated.POST(newContact, controllers.CreateContact)
		authenticated.POST(getContact, controllers.GetContact)
		authenticated.GET(getContacts, controllers.GetContactsFor)
		authenticated.GET(searchContacts, controllers.SearchContactsFor)
		authenticated.POST(getName, controllers.GetName)
		authenticated.GET(getNames, controllers.GetNamesFor)
		authenticated.POST(getAddress, controllers.GetAddress)
		authenticated.POST(getPhone, controllers.GetPhone)
		authenticated.GET(getAddresses, controllers.GetAddressesFor)
		authenticated.GET(getPhones, controllers.GetPhonesFor)
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
