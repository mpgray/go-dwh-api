package app

import (
	u "go-dwh-api/utils"
	"net/http"
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

// UnauthorizedError sends a custom string message and a 403 Error when authentication determines the user is unauthorized
func UnauthorizedError(w http.ResponseWriter, message string) {
	response := make(map[string]interface{})
	response = u.Message(false, message)
	u.Log.Warn(message)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, response)
}
