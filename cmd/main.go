package main

import (
	"context"
	"log"
	"net/http"
	"time"

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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := util.MongoConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	dbName, err := util.GetDatabaseName()
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dbName)

	//TODO: migrate data

	port := util.GetPort()
	router := mux.NewRouter().StrictSlash(true)
	router.Use(handler.LoggingMiddleware)

	router.HandleFunc("/api/health", handler.Health).Methods("GET")
	router.Handle("/api/artwork", handler.GetArtworks(db)).Methods("GET")
	router.Handle("/api/artwork/{id}", handler.GetArtwork(db)).Methods("GET")
	router.Handle("/api/artist", handler.GetArtists(db)).Methods("GET")
	router.Handle("/api/artist/{id}/artwork", handler.GetArtistArtworks(db)).Methods("GET")
	router.Handle("/api/exhibition", handler.GetExhibitions(db)).Methods("GET")
	router.Handle("/api/exhibition/{id}/artwork", handler.GetExhibitionArtworks(db)).Methods("GET")
	router.Handle("/api/exhibition/{id}/artist", handler.GetExhibitionArtists(db)).Methods("GET")

	log.Println("API Started. Listening on", port)
	log.Fatal(http.ListenAndServe(port, router))
}
