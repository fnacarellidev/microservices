package jwtaux

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GetToken(tok string) (*jwt.Token, error) {
	secret, err := os.ReadFile("hs256secret.txt")
	if err != nil {
		return nil, fmt.Errorf("[getToken] read file err %v\n", err)
	}

	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
