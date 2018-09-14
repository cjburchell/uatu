package middelware

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// MySigningKey secret key for JWT tokens
var MySigningKey = "What is the answer"

var middleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return MySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// ValidateJWT token middleware
func ValidateJWT(handlerFunc http.HandlerFunc) http.Handler {
	return middleware.Handler(http.Handler(handlerFunc))
}
