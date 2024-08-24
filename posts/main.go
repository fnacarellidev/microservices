package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func getToken(tok string) (*jwt.Token, error) {
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

func GetPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		log.Println("[GetPosts] Failed at r.Cookie:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	token, err := getToken(cookie.Value)
	if err != nil {
		log.Println("[GetPosts] Failed at r.Cookie:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	tokenMap, _ := token.Claims.(jwt.MapClaims)
	log.Println(tokenMap["username"])
}

func main() {
	router := httprouter.New()
	router.POST("/posts/mine", GetPosts)
	log.Println("Running on port 8081")
	http.ListenAndServe(":8081", router)
}
