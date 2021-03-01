package controllers

import (
	"go-dwh-api/app"
	"go-dwh-api/models"
	u "go-dwh-api/utils"

	"github.com/gin-gonic/gin"
)

// GetStatement gets the current statmen of a single user
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
	resp := u.Message(true, "Contact Retrieved Successfully")
	resp["statement"] = statement
	u.Respond(c.Writer, resp)
}
