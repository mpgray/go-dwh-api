package controllers

import (
	"go-dwh-api/app"
	"go-dwh-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// Login is the controller that manages fetching and comparing email and password
// from the database. It uses bcrypt for the password and returns an access token
// and a refresh token
var Login = func(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		app.UnprocessableEntityError(c, "Invalid json provided during login ")
		return
	}

	dbUser := &models.User{}
	err := app.GetDB().Table("users").Where("email = ?", user.Email).First(dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			app.UnauthorizedError(c, "Log in failed because Email address not found")
			return
		}
		app.InternalServerError(c, "Database Connection error during loggin. Please retry")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		app.ForbiddenError(c, "Invalid login credentials. Please try again")
		return
	}
	//Worked! Logged In
	user.Password = ""

	ts, err := models.CreateToken(dbUser.ID)
	if err != nil {
		app.UnprocessableEntityError(c, "JSON incorrect during the creation of the Tokens")
		return
	}
	saveErr := createAuth(dbUser.ID, ts)
	if saveErr != nil {
		app.UnprocessableEntityError(c, "JSON incorrect during the creation if redis pairs")
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

// Logout removes the redis data from the server.
// by calling DeleteAuth
var Logout = func (c *gin.Context) {
	au, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		app.UnauthorizedError(c, err.Error())
		return
	}
	deleted, delErr := deleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		app.UnauthorizedError(c, "Error logging out")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

// CreateAuth adds the access token and the refresh token to the redis database
func createAuth(userid uint32, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := app.GetRedis().Set(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := app.GetRedis().Set(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func deleteAuth(givenUUID string) (int64, error) {
	deleted, err := app.GetRedis().Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
