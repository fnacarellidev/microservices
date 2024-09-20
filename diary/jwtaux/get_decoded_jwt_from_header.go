package jwtaux

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetDecodedJwtFromCookieHeader(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return nil, err
	}

	token, err := GetToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	tokenMap, _ := token.Claims.(jwt.MapClaims)
	return tokenMap, nil
}
