package jwtaux

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetDecodedJwtFromCookieHeader(jwtCookie http.Cookie) (jwt.MapClaims, error) {
	token, err := GetToken(jwtCookie.Value)
	if err != nil {
		return nil, err
	}

	tokenMap, _ := token.Claims.(jwt.MapClaims)
	return tokenMap, nil
}
