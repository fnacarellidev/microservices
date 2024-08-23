package main

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
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func generateToken(user User) (interface{}, error) {
	secret, err := os.ReadFile("auth/hs256secret.txt")
	if err != nil {
		return nil, fmt.Errorf("[generateToken] read file err %v\n", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("[generateToken] sign token err %v\n", err)
	}

	return tokenString, nil
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user User
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgresql://postgres:postgres@localhost:5432")
	if err != nil {
		log.Println("Failed at db connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close(ctx)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed while parsing body bytes:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		log.Println("Failed to unmarshal body bytes:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := pgquery.New(conn)
	_, err = query.CreateUser(ctx, pgquery.CreateUserParams{
		Username: pgtype.Text{
			String: user.Username,
			Valid: true,
		},
		Password: pgtype.Text{
			String: string(hashedPassword),
			Valid: true,
		},
	})
	if err != nil {
		log.Println("Failed to create user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	router := httprouter.New()
	router.GET("/register", Register)
	http.ListenAndServe(":8080", router)
}
