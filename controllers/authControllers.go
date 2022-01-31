package controllers

import (
	"encoding/json"
	"go-dwh-api/models"
	u "go-dwh-api/utils"
	"net/http"
)

//CreateAccount is a controller to make a new account
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Log.Error("Invalid JSON data recieved when trying to Make a new Account")
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

// Authenticate is a controller for using using
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Log.Error("Invalid JSON data recieved when trying to Authenticate an Account")
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
