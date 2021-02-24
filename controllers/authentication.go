package controllers

import (
	"go-dwh-api/app"
	"go-dwh-api/models"
	u "go-dwh-api/utils"
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
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	dbUser := &models.User{}
	err := app.GetDB().Table("users").Where("email = ?", user.Email).First(dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warn("Attempt to login resulted in email not found.")
			c.JSON(http.StatusForbidden, u.Message(false, "Login Failed: Email address not found"))
			return
		}
		u.Log.Error("DB Connection Failed: During login attempt")
		c.JSON(http.StatusForbidden, u.Message(false, "Connection error. Please retry"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		u.Log.Warn("Attempt to login with invalid credentials: %s ", err.Error())
		c.JSON(http.StatusForbidden, u.Message(false, "Invalid login credentials. Please try again"))
		return
	}
	//Worked! Logged In
	user.Password = ""

	ts, err := models.CreateToken(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := createAuth(user.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
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
func Logout(c *gin.Context) {
	au, err := models.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := deleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

// CreateAuth adds the access token and the refresh token to the redis database
func createAuth(userid uint, td *models.TokenDetails) error {
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
