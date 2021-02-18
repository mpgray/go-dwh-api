package controllers

import (
	"encoding/json"
	"go-dwh-api/models"
	u "go-dwh-api/utils"
	"net/http"
)

// CreateContact is a controller to make a new contact
var CreateContact = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		u.Log.Error("Invalid JSON data recieved when trying to make a new Contact.")
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	contact.UserID = user
	resp := contact.Create()
	u.Respond(w, resp)
}

// GetContact Gets all the contact information for a single user.
var GetContact = func(w http.ResponseWriter, r *http.Request) {
	contactID := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contactID)
	if err != nil {
		u.Log.Error("Invalid JSON data recieved when getting a contact.")
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	userID := r.Context().Value("user").(uint)
	data := models.GetContact(contactID.ID, userID)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

// GetContactsFor gets all the contacts associated with an owner
var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user").(uint)
	data := models.GetContacts(userID)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
