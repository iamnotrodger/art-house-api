package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/internal/handler"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := util.GetPort()
	router := mux.NewRouter().StrictSlash(true)
	router.Use(handler.LoggingMiddleware)

	router.HandleFunc("/api/health", handler.Health).Methods("GET")

	log.Println("API Started. Listening on", port)
	log.Fatal(http.ListenAndServe(port, router))
}
