package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fnacarellidev/microsservices/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/microsservices/auth/api"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user api.User
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
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
		Username: user.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		log.Println("Failed to create user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
