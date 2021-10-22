package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/internal/artist"
	"github.com/iamnotrodger/art-house-api/internal/artwork"
	"github.com/iamnotrodger/art-house-api/internal/handler"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()

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

	artworkHandler := artwork.NewHandler(db)
	artistHandler := artist.NewHandler(db)

	port := util.GetPort()
	router := mux.NewRouter().StrictSlash(true)
	router.Use(handler.LoggingMiddleware)

	//Health Routes
	router.HandleFunc("/api/health", handler.Health).Methods("GET")

	//Artwork Routes
	router.HandleFunc("/api/artwork", artworkHandler.GetMany).Methods("GET")
	router.HandleFunc("/api/artwork/{id}", artworkHandler.Get).Methods("GET")

	//Artist Routes
	router.HandleFunc("/api/artist", artistHandler.GetMany).Methods("GET")
	router.HandleFunc("/api/artist/{id}", artistHandler.Get).Methods("GET")
	router.HandleFunc("/api/artist/{id}/artwork", artistHandler.GetArtworks).Methods("GET")

	//Exhibition Routes
	router.Handle("/api/exhibition", handler.GetExhibitions(db)).Methods("GET")
	router.Handle("/api/exhibition/{id}", handler.GetExhibition(db)).Methods("GET")
	router.Handle("/api/exhibition/{id}/artwork", handler.GetExhibitionArtworks(db)).Methods("GET")
	router.Handle("/api/exhibition/{id}/artist", handler.GetExhibitionArtists(db)).Methods("GET")

	server := cors.Default().Handler(router)
	log.Println("API Started. Listening on", port)
	log.Fatal(http.ListenAndServe(port, server))
}
