package app

import (
	u "go-dwh-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler Error message when http message not found
var NotFoundHandler = func(next http.Handler) http.Handler {
	message := "This resource was not found on our server"
	u.Log.Error(message)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, message))
		next.ServeHTTP(w, r)
	})
}

// UnauthorizedError sends a custom string message and a 401 Error when authentication determines the user is unauthorized
// Unauthorized means that the user could be allowed but doesn't have proper credentials
func UnauthorizedError(c *gin.Context, message string) {
	response := make(map[string]interface{})
	response = u.Message(false, message)
	u.Log.Info(message)
	c.JSON(http.StatusUnauthorized, response)
}

// ForbiddenError sends a custom string message and a 403 error
// Forbidden mean the user is not allowed to view this under any circumtaces.
func ForbiddenError(c *gin.Context, message string) {
	response := make(map[string]interface{})
	response = u.Message(false, message)
	u.Log.Warn(message)
	c.JSON(http.StatusForbidden, response)
}

// InternalServerError sends a cusomt string message and a 500 error
// This is a generic error that just means to the user that it is server
// side and nothing to do with them or their client
func InternalServerError(c *gin.Context, message string) {
	response := make(map[string]interface{})
	response = u.Message(false, message)
	u.Log.Error(message)
	c.JSON(http.StatusInternalServerError, response)
}

// UnprocessableEntityError sends a custom string message and 422 error
// This usually means that we recieved bad json data from the client
func UnprocessableEntityError(c *gin.Context, message string) {
	response := make(map[string]interface{})
	response = u.Message(false, message)
	u.Log.Error(message)
	c.JSON(http.StatusUnprocessableEntity, response)
}
