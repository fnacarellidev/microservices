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

func GetRecords(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Println("[GetRecords] Failed at pgx.Connect:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close(ctx)
	queries := pgquery.New(conn)
	decodedJwt, err := jwtaux.GetDecodedJwtFromCookieHeader(r)
	if err != nil {
		log.Println("[GetRecords] Failed at jwtaux.GetDecodedJwtFromCookieHeader:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	records, err := queries.GetRecordsFromUser(ctx, decodedJwt["username"].(string))
	if err != nil {
		log.Println("[GetRecords] Failed at queries.GetPostsFromUser:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, record := range records {
		log.Println(record)
	}
}
