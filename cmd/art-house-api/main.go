package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/cmd/config"
	"github.com/iamnotrodger/art-house-api/internal/artist"
	"github.com/iamnotrodger/art-house-api/internal/artwork"
	"github.com/iamnotrodger/art-house-api/internal/exhibition"
	"github.com/iamnotrodger/art-house-api/internal/health"
	"github.com/iamnotrodger/art-house-api/internal/middleware"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"github.com/rs/cors"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := util.GetMongoClient(ctx, config.Global.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database(config.Global.MongoDBName)

	// TODO: migrate data

	artworkHandler := artwork.NewHandler(db)
	artistHandler := artist.NewHandler(db)
	exhibitionHandler := exhibition.NewHandler(db)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.LoggingMiddleware)

	//Health Routes
	router.HandleFunc("/api/health", health.GetHealth).Methods("GET")
	//Artwork Routes
	artworkHandler.RegisterRoutes(router)
	//Artist Routes
	artistHandler.RegisterRoutes(router)
	//Exhibition Routes
	exhibitionHandler.RegisterRoutes(router)

	server := cors.Default().Handler(router)
	log.Println("API Started. Listening on", config.Global.Port)
	log.Fatal(http.ListenAndServe(config.Global.Port, server))
}
