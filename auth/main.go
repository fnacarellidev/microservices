package main

import (
	"log"
	"net/http"

	"github.com/fnacarellidev/microsservices/auth/handlers"
	"github.com/julienschmidt/httprouter"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := httprouter.New()
	router.POST("/auth/register", handlers.RegisterHandler)
	router.POST("/auth/login", handlers.LoginHandler)
	log.Println("Running on port 8080")
	http.ListenAndServe(":8080", CORS(router))
}
