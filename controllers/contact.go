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
	resp["contacts"] = contacts
	u.Respond(c.Writer, resp)
}

// GetNamesFor gets all the names associated with an owner
var GetNamesFor = func(c *gin.Context) {
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get names ")
		return
	}
	userID := metadata.UserID

	names := models.GetAllContactInfo(userID, models.NAME)
	resp := u.Message(true, "All names retrieved successfully.")
	resp["names"] = names
	u.Respond(c.Writer, resp)
}

var GetName = func(c *gin.Context) {
	contactID := &models.ContactID{}
	userID, err := models.FetchAuthenticatedID(c, &contactID)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get Contact "+models.NAME)
		return
	}

	name := models.GetContactInfo(userID, contactID.ID, string(models.NAME))
	if name == nil {
		app.ForbiddenError(c, "That user isn't associated with you.")
		return
	}
	resp := u.Message(true, models.NAME+" Retrieved Successfully")
	resp[models.NAME] = name
	u.Respond(c.Writer, resp)
}

var GetAddress = func(c *gin.Context) {
	contactID := &models.ContactID{}
	userID, err := models.FetchAuthenticatedID(c, &contactID)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get Contact "+models.ADDRESS)
		return
	}

	address := models.GetContactInfo(userID, contactID.ID, string(models.ADDRESS))
	if address == nil {
		app.ForbiddenError(c, "That user isn't associated with you.")
		return
	}
	resp := u.Message(true, models.ADDRESS+" Retrieved Successfully")
	resp[models.NAME] = address
	u.Respond(c.Writer, resp)
}

var GetAddressesFor = func(c *gin.Context) {
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get addresses ")
		return
	}
	userID := metadata.UserID

	addresses := models.GetAllContactInfo(userID, models.ADDRESS)
	resp := u.Message(true, "All Addresses retrieved successfully.")
	resp["addresses"] = addresses
	u.Respond(c.Writer, resp)
}

var GetPhone = func(c *gin.Context) {
	contactID := &models.ContactID{}
	userID, err := models.FetchAuthenticatedID(c, &contactID)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get Contact "+models.PHONE)
		return
	}

	phone := models.GetContactInfo(userID, contactID.ID, string(models.PHONE))
	if phone == nil {
		app.ForbiddenError(c, "That user isn't associated with you.")
		return
	}
	resp := u.Message(true, models.PHONE+" Retrieved Successfully")
	resp[models.PHONE] = phone
	u.Respond(c.Writer, resp)
}

var GetPhonesFor = func(c *gin.Context) {
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get phone numbers ")
		return
	}
	userID := metadata.UserID

	phones := models.GetAllContactInfo(userID, models.PHONE)
	resp := u.Message(true, "All Phone Numbers retrieved successfully.")
	resp["phones"] = phones
	u.Respond(c.Writer, resp)
}

var SearchContactsFor = func(c *gin.Context) {
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get contact search ")
		return
	}
	userID := metadata.UserID

	phones := models.SearchContacts(userID)
	resp := u.Message(true, "Searching...")
	resp["search"] = phones
	u.Respond(c.Writer, resp)
}
