package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fnacarellidev/microsservices/.sqlcbuild/pgquery"
	"github.com/fnacarellidev/microsservices/posts/api"
	"github.com/fnacarellidev/microsservices/posts/jwtaux"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var post api.Post
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Println("[CreatePost] Failed at pgx.Connect:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close(ctx)
	queries := pgquery.New(conn)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[GetPosts] Failed at io.ReadAll:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := json.Unmarshal(bodyBytes, &post); err != nil {
		log.Println("[GetPosts] Failed at json.Unmarshal:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	decodedJwt, err := jwtaux.GetDecodedJwtFromCookieHeader(r)
	if err != nil {
		log.Println("[CreatePost] Failed at jwtaux.GetDecodedJwtFromCookieHeader:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	id, err := queries.CreatePost(ctx, pgquery.CreatePostParams{
		PostOwner: decodedJwt["username"].(string),
		Content: post.Content,
	})
	if err != nil {
		log.Println("[GetPosts] Failed at queries.CreatePost:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	res := struct {
		Id string `json:"id"`
	}{
		Id: id.String(),
	}
	body, err := json.Marshal(res)
	if err != nil {
		log.Println("[GetPosts] Failed at json.Marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
