package middelware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/auth0/go-jwt-middleware"
	"net/http"
)

var MySigningKey = "What is the answer"

var middleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return MySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func ValidateJWT(handlerFunc http.HandlerFunc) http.Handler{
	return middleware.Handler(http.Handler(handlerFunc))
}
