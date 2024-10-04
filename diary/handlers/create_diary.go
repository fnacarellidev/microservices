package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/fnacarellidev/microsservices/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/microsservices/diary/jwtaux"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func CreateDiary(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Println("[GetRecords] Failed at pgx.Connect:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close(ctx)
	queries := pgquery.New(conn)
	jwt, err := r.Cookie("jwt")
	if err != nil {
		log.Println("[GetRecords] Failed at jwtaux.GetDecodedJwtFromCookieHeader:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	decodedJwt, err := jwtaux.GetDecodedJwtFromCookieHeader(*jwt)
	if err != nil {
		log.Println("[GetRecords] Failed at jwtaux.GetDecodedJwtFromCookieHeader:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userId, err := queries.GetIdFromUser(ctx, decodedJwt["username"].(string))
	if err != nil {
		log.Println("[GetRecords] Failed at GetIdFromUser:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := queries.CreateDiaryForUser(ctx, userId); err != nil {
		log.Println("[GetRecords] Failed at CreateDiaryForUser:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
}
