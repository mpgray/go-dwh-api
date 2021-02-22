package main

import (
	"fmt"
	"go-dwh-api/api"
	"go-dwh-api/app"
	m "go-dwh-api/models"
	u "go-dwh-api/utils"
	"net/http"
	"os"
)

func main() {
	app.GetDB().AutoMigrate(&m.Account{}, &m.Contact{}, &m.FullName{},
		&m.Address{}, &m.Phone{}, &m.Statement{}, &m.User{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8989" //localhost
	} else {
		u.Log.Fatalf("Attempted to connect to port %s but it was already in use.", port)
	}

	u.Log.Info("Connected on port " + port)

	router := api.Router()
	handler := api.CorsConfig().Handler(router)
	err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:8989/api
	if err != nil {
		u.Log.Fatal(fmt.Sprint(err))
	}

	defer u.Log.Infof("**Golang Backend API for Driveway Home Started Successfully**")
}
