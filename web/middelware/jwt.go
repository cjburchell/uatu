package middelware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/auth0/go-jwt-middleware"
)

// MySigningKey secret key for JWT tokens
var MySigningKey = "What is the answer"

var middleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(MySigningKey), nil
	},
	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodHS256,
})

// ValidateJWT token middleware
func ValidateJWT(handlerFunc http.HandlerFunc) http.Handler {
	return middleware.Handler(http.Handler(handlerFunc))
}
