package controllers

import (
	"go-dwh-api/app"
	"go-dwh-api/models"
	u "go-dwh-api/utils"

	"github.com/gin-gonic/gin"
)

//CreateUser is a controller to make a new account
var CreateUser = func(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		app.UnprocessableEntityError(c, "Invalid json provided during Creation of a user account "+err.Error())
		return
	}

	resp := user.CreateUser() //Create account
	u.Respond(c.Writer, resp)
}
