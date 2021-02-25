package controllers

import (
	"go-dwh-api/app"
	"go-dwh-api/models"
	u "go-dwh-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateContact is a controller to make a new contact
var CreateContact = func(c *gin.Context) {
	contact := &models.Contact{}
	userID, err := models.FetchAuthenticatedID(c, &contact)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get a contact.")
		return
	}

	contact.UserID = userID
	resp := contact.CreateContact() //Create contact
	u.Respond(c.Writer, resp)
}

// GetContact Gets all the contact information for a single user.
var GetContact = func(c *gin.Context) {

	contactID := &models.ContactID{}
	userID, err := models.FetchAuthenticatedID(c, &contactID)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get a contact ")
		return
	}

	contact := models.GetContact(contactID.ID, userID)
	if contact == nil {
		app.ForbiddenError(c, "That user isn't associated with you.")
		return
	}
	resp := u.Message(true, "Contact Retrieved Successfully")
	resp["contact"] = contact
	u.Respond(c.Writer, resp)
}

// GetContactsFor gets all the contacts associated with an owner
var GetContactsFor = func(c *gin.Context) {
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get contacts ")
		return
	}
	userID := metadata.UserID

	contacts := models.GetContacts(userID)
	resp := u.Message(true, "All Contacts retrieved successfully.")
	resp["data"] = contacts
	u.Respond(c.Writer, resp)
}
