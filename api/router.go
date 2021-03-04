package api

import (
	"go-dwh-api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

// Router creates and serves the server
func Router() *gin.Engine {
	router := gin.Default()
	const API string = "/api/v1"

	// Group that needs no authenticating, i.e. unauthenticated
	unauthenticated := router.Group(API)
	{
		unauthenticated.POST("/login", controllers.Login)
		unauthenticated.POST("/user/new", controllers.CreateUser)

	}

	// Group that requires an authenticated user
	authenticated := router.Group(API + "/auth")
	authenticated.Use(controllers.TokenAuthenticator())
	{
		authenticated.POST("/refresh", controllers.Refresh)
		authenticated.POST("/logout", controllers.Logout)
	}

	contacts := router.Group(API + "/contact")
	contacts.Use(controllers.TokenAuthenticator())
	{
		contacts.POST("", controllers.GetContact)
		contacts.POST("/new", controllers.CreateContact)
		contacts.GET("/all", controllers.GetContactsFor)
		contacts.GET("/search", controllers.SearchContactsFor)
		contacts.POST("/name", controllers.GetName)
		contacts.GET("/names", controllers.GetNamesFor)
		contacts.POST("/address", controllers.GetAddress)
		contacts.POST("/phone", controllers.GetPhone)
		contacts.GET("/addresses", controllers.GetAddressesFor)
		contacts.GET("/phones", controllers.GetPhonesFor)
	}

	statements := router.Group(API + "/statement")
	statements.Use(controllers.TokenAuthenticator())
	{
		statements.POST("", controllers.GetStatement)
		statements.POST("/new", controllers.CreateStatement)
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
