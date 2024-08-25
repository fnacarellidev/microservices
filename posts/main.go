package main

import (
	"log"
	"net/http"

	"github.com/fnacarellidev/microsservices/posts/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/posts/mine", handlers.GetPostsHandler)
	router.POST("/posts/create", handlers.CreatePostHandler)
	router.DELETE("/posts/delete/:id", handlers.DeletePostHandler)
	log.Println("Running on port 8081")
	http.ListenAndServe(":8081", router)
}
