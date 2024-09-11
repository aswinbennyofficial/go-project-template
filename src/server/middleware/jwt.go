package middlewares

import (
	"github.com/go-chi/jwtauth/v5"
)

func InitJWTAuth(secret string) (*jwtauth.JWTAuth){
	return jwtauth.New("HS256", []byte(secret), nil)
}


