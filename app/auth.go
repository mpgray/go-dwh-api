package app

import (
	"context"
	"fmt"
	"go-hoa-api/models"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// JwtAuthentication resolves if the token is valid, exists and matches.
// If it does not, it returns a 403 error.
var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/v1/user/new", "/v1/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                             //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			UnauthorizedError(w, "Missing auth token")
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			UnauthorizedError(w, "Invalid/Malformed auth token")
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			UnauthorizedError(w, "Malformed authentication token")
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			UnauthorizedError(w, "Token is not valid.")
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		// fmt.Sprintf("User %d", tk.UserID) //Useful for monitoring
		fmt.Printf("User %d", tk.UserID)
		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
