package artist

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/query"
	"github.com/iamnotrodger/art-house-api/internal/util"
)

type Handler struct {
	store *Store
	cache *Cache
}

func NewHandler(store *Store, cache *Cache) *Handler {
	return &Handler{
		store: store,
		cache: cache,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/artist", h.GetMany).Methods("GET")
	router.HandleFunc("/api/artist/{id}", h.Get).Methods("GET")
	router.HandleFunc("/api/artist/{id}/artwork", h.GetArtworks).Methods("GET")
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	artistID := params["id"]

	artist, err := h.getOrSetArtistCache(r.Context(), artistID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(artist)
}

func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryString := r.URL.RawQuery
	artists, err := h.cache.GetMany(r.Context(), queryString)
	if err != nil {
		log.Println(err)
	} else if artists != nil {
		json.NewEncoder(w).Encode(artists)
		return
	}

	queryParams := query.NewArtistQuery(r.URL.Query())
	artists, err = h.store.FindMany(r.Context(), queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	err = h.cache.SetMany(r.Context(), queryString, artists)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(artists)
}

func (h *Handler) GetArtworks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	artistID := params["id"]

	queryString := r.URL.RawQuery
	artworks, err := h.cache.GetArtworks(r.Context(), artistID, queryString)
	if err != nil {
		log.Println(err)
	} else if artworks != nil {
		json.NewEncoder(w).Encode(artworks)
		return
	}

	queryParams := query.NewArtworkQuery(r.URL.Query())
	artworks, err = h.store.FindArtworks(r.Context(), artistID, queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	artist, err := h.getOrSetArtistCache(r.Context(), artistID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	for _, artwork := range artworks {
		artwork.Artist = artist
	}

	err = h.cache.SetArtworks(r.Context(), artistID, queryString, artworks)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(artworks)
}

func (h *Handler) getOrSetArtistCache(ctx context.Context, artistID string) (*model.Artist, error) {
	artist, err := h.cache.Get(ctx, artistID)
	if err != nil {
		log.Println(err)
	} else if artist != nil {
		return artist, nil
	}

	artist, err = h.store.Find(ctx, artistID)
	if err != nil {
		return nil, err
	}
	err = h.cache.Set(ctx, artistID, artist)
	if err != nil {
		log.Println(err)
	}

	return artist, nil
}
