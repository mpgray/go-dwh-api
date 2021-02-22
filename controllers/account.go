package controllers

import (
	"go-dwh-api/models"
	u "go-dwh-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//CreateAccount is a controller to make a new account
var CreateAccount = func(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	resp := user.CreateUser() //Create account
	u.Respond(c.Writer, resp)
}
