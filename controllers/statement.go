package controllers

import (
	"go-dwh-api/app"
	"go-dwh-api/models"
	u "go-dwh-api/utils"

	"github.com/gin-gonic/gin"
)

var GetStatement = func(c *gin.Context) {
	contactID := &models.ContactID{}
	_, err := models.FetchAuthenticatedID(c, &contactID)
	if err != nil {
		app.UnauthorizedError(c, "Unauthorized attempt to get a statement ")
		return
	}

	statement := models.GetCurrentStatement(contactID.ID)
	if contact == nil {
		app.ForbiddenError(c, "That user isn't associated with you.")
		return
	}
	resp := u.Message(true, "Contact Retrieved Successfully")
	resp["statement"] = contact
	u.Respond(c.Writer, resp)
}
