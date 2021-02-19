package main

import (
	"fmt"
	"go-dwh-api/app"
	u "go-dwh-api/utils"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8989" //localhost
	} else {
		u.Log.Fatalf("Attempted to connect to port %s but it was already in use.", port)
	}

	u.Log.Info("Connected on port " + port)

	router := app.Router()
	handler := app.CorsConfig().Handler(router)
	err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:8989/api
	if err != nil {
		u.Log.Fatal(fmt.Sprint(err))
	}

	defer u.Log.Infof("**Golang Backend API for Driveway Home Started Successfully**")
}
