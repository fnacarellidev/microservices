package handlers

import (
	"context"
	"encoding/json"
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

func GetPostsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var getPostsRes api.GetPostsRes
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Println("[GetPosts] Failed at pgx.Connect:", err)
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

	tokenMap, _ := token.Claims.(jwt.MapClaims)
	posts, err := queries.GetPostsFromUser(ctx, tokenMap["username"].(string))
	if err != nil {
		log.Println("[GetPosts] Failed at queries.GetPostsFromUser:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, post := range posts {
		getPostsRes.Posts = append(getPostsRes.Posts, api.Post{
			Content: post.Content,
			CreatedAt: post.CreatedAt.Time,
		})
	}

	body, err := json.Marshal(getPostsRes)
	if err != nil {
		log.Println("[GetPosts] Failed at json.Marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
