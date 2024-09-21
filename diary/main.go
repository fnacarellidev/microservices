package main

import (
	"log"
	"net/http"

	"github.com/fnacarellidev/microsservices/diary/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/diary/my_records", handlers.GetRecords)
	router.POST("/diary/create_record", handlers.CreateRecord)
	router.POST("/diary/create_diary", handlers.CreateDiary)
	log.Println("Running on port 8081")
	http.ListenAndServe(":8081", router)
}
