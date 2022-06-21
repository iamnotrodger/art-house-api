package main

import (
	"context"
	"fmt"
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
	config.LoadConfig()

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := util.GetMongoClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	db := client.Database(config.Global.MongoDBName)

	rdb := util.GetRedisClient()

	// TODO: create app context to hold all the db and cache
	artworkStore := artwork.NewStore(db)
	artworkCache := artwork.NewCache(rdb, time.Minute)
	artworkHandler := artwork.NewHandler(artworkStore, artworkCache)

	artistStore := artist.NewStore(db)
	artistCache := artist.NewCache(rdb, time.Minute)
	artistHandler := artist.NewHandler(artistStore, artistCache)

	exhibitionStore := exhibition.NewStore(db)
	exhibitionCache := exhibition.NewCache(rdb, time.Minute)
	exhibitionHandler := exhibition.NewHandler(exhibitionStore, exhibitionCache)

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

	server := &http.Server{
		Handler:      cors.Default().Handler(router),
		Addr:         fmt.Sprintf("0.0.0.0:%v", config.Global.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("API Started. Listening on", config.Global.Port)
	log.Fatal(server.ListenAndServe())
}
