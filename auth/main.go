package main

import (
	"log"
	"net/http"

	"github.com/fnacarellidev/microsservices/auth/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.POST("/auth/register", handlers.RegisterHandler)
	router.POST("/auth/login", handlers.LoginHandler)
	log.Println("Running on port 8080")
	http.ListenAndServe(":8080", router)
}
