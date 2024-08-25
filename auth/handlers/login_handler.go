package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fnacarellidev/microsservices/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/microsservices/auth/api"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func generateToken(user api.User) (string, error) {
	secret, _ := os.ReadFile("hs256secret.txt")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("[generateToken] sign token err %v\n", err)
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user api.User
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Println("[LOGIN] pgx.Connect:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close(ctx)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[LOGIN] io.ReadAll:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		log.Println("[LOGIN] json.Unmarshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	queries := pgquery.New(conn)
	hashedPassword, err := queries.GetPasswordFromUser(ctx, user.Username)
	if err != nil {
		log.Println("[LOGIN] queries.GetPasswordFromUser:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		log.Println("[LOGIN] bcrypt.CompareHashAndPassword:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := generateToken(user)
	if err != nil {
		log.Println("[LOGIN] jwtaux.GenerateToken:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jwtCookie := http.Cookie{
		Name: "jwt",
		Value: token,
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &jwtCookie)
}
