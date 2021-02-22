package controllers

import (
	"encoding/json"
	"fmt"
	"go-dwh-api/models"
	u "go-dwh-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateContact is a controller to make a new contact
var CreateContact = func(c *gin.Context) {

	user := c.Query("user") //Grab the id of the user that send the request
	fmt.Printf("User that created contact %s", user)
	contact := &models.Contact{}

	if err := c.ShouldBindJSON(&contact); err != nil {
		u.Log.Error("Invalid JSON data recieved when trying to make a new Contact.")
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	resp := contact.CreateContact() //Create account
	u.Respond(c.Writer, resp)
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
