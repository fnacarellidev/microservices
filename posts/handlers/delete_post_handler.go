package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/fnacarellidev/microsservices/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/microsservices/posts/jwtaux"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idParam := p.ByName("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Println("[DeletePostHandler] Failed at uuid.Parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Println("[DeletePostsHandler] Failed at pgx.Connect:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close(ctx)
	queries := pgquery.New(conn)
	decodedJwt, err := jwtaux.GetDecodedJwtFromCookieHeader(r)
	if err != nil {
		log.Println("[DeletePostsHandler] Failed at jwtaux.GetDecodedJwtFromCookieHeader:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	username := decodedJwt["username"].(string)
	posts, err := queries.GetPostsFromUser(ctx, username)
	if err != nil {
		log.Println("[DeletePostsHandler] Failed at queries.GetPostsFromUser:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, post := range posts {
		if post.ID == id {
			if err := queries.DeletePostById(ctx, id); err != nil {
				log.Println("[DeletePostsHandler] Failed at queries.DeletePostById:", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	return
}
