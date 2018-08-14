package login

import (
	"net/http"
	"crypto/subtle"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/cjburchell/yasls/web/middelware"
)

func SetupRoutes(r *mux.Router, username string, password string)  {
	r.HandleFunc("/login", basicAuthLogin(username, password)).Methods("POST")
}

func basicAuthLogin(username string, password string) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		/* Create the token */
		token := jwt.New(jwt.SigningMethodHS256)

		// Create a map to store our claims
		claims := token.Claims.(jwt.MapClaims)

		/* Set token claims */
		claims["name"] = username
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		/* Sign the token with our secret */
		tokenString, _ := token.SignedString(middelware.MySigningKey)

		/* Finally, write the token to the browser window */
		w.Write([]byte(tokenString))
	}
}
