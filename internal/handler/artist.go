package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetArtist(db *mongo.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{ "message": "Invalid ID"}`))
			return
		}

		artist, err := model.FindArtist(db, id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
			return
		}

		json.NewEncoder(w).Encode(artist)
	})
}

func GetArtists(db *mongo.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		options := util.QueryBuilder(r.URL.Query())
		artists, err := model.FindArtists(db, bson.D{}, options)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
			return
		}

		json.NewEncoder(w).Encode(artists)
	})
}

func GetArtistArtworks(db *mongo.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{ "message": "Invalid ID"}`))
			return
		}

		queryParams := r.URL.Query()
		delete(queryParams, "search")

		options := util.QueryBuilderPipeline(queryParams)
		artworks, err := model.FindArtistArtworks(db, id, options...)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
			return
		}

		json.NewEncoder(w).Encode(artworks)
	})
}