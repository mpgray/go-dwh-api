package controllers

import (
	"go-dwh-api/app"
	"go-dwh-api/models"
	u "go-dwh-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateStatement is a controller to make a new contact
var CreateStatement = func(c *gin.Context) {
	contact := &models.Contact{}
	userID, err := models.FetchAuthenticatedID(c, &contact)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to create a statement.")
		return
	}

	contact.UserID = userID
	statement := &models.Statement{}
	resp := statement.CreateStatement() //Create contact
	u.Respond(c.Writer, resp)
}

// GetStatement gets the current statment of a single user
var GetStatement = func(c *gin.Context) {
	contactID := &models.ContactID{}
	_, err := models.FetchAuthenticatedID(c, &contactID)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get a statement ")
		return
	}

	statement := models.GetCurrentStatement(contactID.ID)
	if statement == nil {
		app.ForbiddenError(c, "That statement isn't associated with you.")
		return
	}
	resp := u.Message(true, "Statment Retrieved Successfully")
	resp["statement"] = statement
	u.Respond(c.Writer, resp)
}

// GetStatementsFor searches for contact by address and name
var GetStatementsFor = func(c *gin.Context) {
	metadata, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get statement history ")
		return
	}
	userID := metadata.UserID

	statements := models.GetStatementHistory(userID)
	resp := u.Message(true, "Statement History Retrieved Successfully")
	resp["statement"] = statements
	u.Respond(c.Writer, resp)
}
