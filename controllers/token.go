package controllers

import (
	"fmt"
	"go-dwh-api/app"
	"go-dwh-api/models"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// TokenAuthenticator checks if the user is autherized
func TokenAuthenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := models.TokenValid(c.Request)
		if err != nil {
			app.UnauthorizedError(c, "You are not authorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Refresh uses the refresh token to get a new Authentication and
// Refresh token.
var Refresh = func(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		app.UnprocessableEntityError(c, "Incorrect Json during refresh")
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		app.UnauthorizedError(c, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		app.UnauthorizedError(c, "Token Invalid")
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			app.UnprocessableEntityError(c, "Bad claim information")
			return
		}
		userID64, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 32)
		userID := uint32(userID64)

		if err != nil {
			app.UnprocessableEntityError(c, "JSON Error during user id from cache.")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := deleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			app.UnauthorizedError(c, "That refresh token does't exist.")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := models.CreateToken(userID)
		if createErr != nil {
			app.ForbiddenError(c, "Restriced access during token creation")
			return
		}
		//save the tokens metadata to redis
		saveErr := createAuth(userID, ts)
		if saveErr != nil {
			app.ForbiddenError(c, "Restriced access during redis cache creation")
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		app.UnauthorizedError(c, "Your session has expired")
	}
}
