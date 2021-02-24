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
	if err := c.ShouldBindJSON(&contact); err != nil {
		app.UnprocessableEntityError(c, "Invalid JSON recieved during Creation of contact "+err.Error())
		return
	}
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, err.Error())
		return
	}
	/*
		userID, err := models.FetchAuth(metadata)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		} */

	contact.UserID = metadata.UserID
	resp := u.Message(true, "Contact created successfully")
	resp["contact"] = contact.CreateContact() //Create contact
	u.Respond(c.Writer, resp)
}

// GetContact Gets all the contact information for a single user.
var GetContact = func(c *gin.Context) {
	/*	contactID := &models.ContactID{}
		if err := c.ShouldBindJSON(&contactID); err != nil {
			app.UnprocessableEntityError(c, "Invalid JSON recieved during getting of a contact "+err.Error())
			return
		}

		metadata, err := models.ExtractTokenMetadata(c.Request)
		if err != nil {
			app.UnauthorizedError(c, "Unauthorized attempt to get a contact "+err.Error())
			return
		} */
	contactID := &models.ContactID{}
	userID, err := models.FetchAuthenticatedID(c, &contactID)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get a contact "+err.Error())
		return
	}

	data := models.GetContact(contactID.ID, userID)
	if data == nil {
		app.ForbiddenError(c, "That user isn't associated with you. "+err.Error())
		return
	}
	resp := u.Message(true, "Contact Retrieved Successfully")
	resp["data"] = data
	u.Respond(c.Writer, resp)
}

// GetContactsFor gets all the contacts associated with an owner
var GetContactsFor = func(c *gin.Context) {
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get contacts "+err.Error())
		return
	}
	userID := metadata.UserID

	contacts := models.GetContacts(userID)
	resp := u.Message(true, "All Contacts retrieved successfully.")
	resp["contacts"] = contacts
	u.Respond(c.Writer, resp)
}
