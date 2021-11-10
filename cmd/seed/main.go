package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/iamnotrodger/art-house-api/internal/artist"
	"github.com/iamnotrodger/art-house-api/internal/artwork"
	"github.com/iamnotrodger/art-house-api/internal/exhibition"
	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"github.com/joho/godotenv"
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

	artworkStore := artwork.NewStore(db)
	artistStore := artist.NewStore(db)
	exhibitionStore := exhibition.NewStore(db)

	artists, err := parseArtists("./cmd/seed/data/artists.json")
	if err != nil {
		log.Fatal(err)
	}

	artworks, err := parseArtworks("./cmd/seed/data/artworks.json")
	if err != nil {
		log.Fatal(err)
	}

	exhibitions, err := parseExhibitions("./cmd/seed/data/exhibitions.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = artistStore.InsertMany(context.TODO(), artists)
	if err != nil {
		log.Fatal(err)
	}

	_, err = artworkStore.InsertMany(context.TODO(), artworks)
	if err != nil {
		log.Fatal(err)
	}

	_, err = exhibitionStore.InsertMany(context.TODO(), exhibitions)
	if err != nil {
		log.Fatal(err)
	}
}

func parseArtists(filePath string) ([]model.Artist, error) {
	var artists []model.Artist

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(file), &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func parseArtworks(filePath string) ([]model.Artwork, error) {
	var artworks []model.Artwork

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(file), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks, nil
}

func parseExhibitions(filePath string) ([]model.Exhibition, error) {
	var exhibitions []model.Exhibition

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(file), &exhibitions)
	if err != nil {
		return nil, err
	}

	return exhibitions, nil
}
