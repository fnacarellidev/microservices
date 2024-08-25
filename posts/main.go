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
	"github.com/julienschmidt/httprouter"
)

type Post struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type GetPostsRes struct {
	Posts []Post `json:"posts"`
}

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
	var getPostsRes GetPostsRes
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

	token, err := getToken(cookie.Value)
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
		getPostsRes.Posts = append(getPostsRes.Posts, Post{
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

func CreatePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var post Post
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

	token, err := getToken(cookie.Value)
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

func main() {
	router := httprouter.New()
	router.GET("/posts/mine", GetPosts)
	router.POST("/posts/create", CreatePost)
	log.Println("Running on port 8081")
	http.ListenAndServe(":8081", router)
}
