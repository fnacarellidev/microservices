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
	"github.com/golang-jwt/jwt/v5"
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
	cookie, err := r.Cookie("jwt")
	if err != nil {
		log.Println("[GetPosts] Failed at r.Cookie:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	token, err := jwtaux.GetToken(cookie.Value)
	if err != nil {
		log.Println("[GetPosts] Failed at r.Cookie:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

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

	tokenMap, _ := token.Claims.(jwt.MapClaims)
	id, err := queries.CreatePost(ctx, pgquery.CreatePostParams{
		PostOwner: tokenMap["username"].(string),
		Content: post.Content,
	})
	if err != nil {
		log.Println("[GetPosts] Failed at queries.CreatePost:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Write([]byte(id.String()))
}
