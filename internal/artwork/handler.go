package artwork

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	router.HandleFunc("/api/artwork", h.GetMany).Methods("GET")
	router.HandleFunc("/api/artwork/{id}", h.Get).Methods("GET")
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	artworkID := params["id"]

	artwork, err := h.cache.Get(r.Context(), artworkID)
	if err != nil {
		log.Println(err)
	} else if artwork != nil {
		json.NewEncoder(w).Encode(artwork)
		return
	}

	artwork, err = h.store.Find(r.Context(), artworkID)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	err = h.cache.Set(r.Context(), artworkID, artwork)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(artwork)
}

func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryString := r.URL.RawQuery

	artworks, err := h.cache.GetMany(r.Context(), queryString)
	if err != nil {
		log.Println(err)
	} else if artworks != nil {
		json.NewEncoder(w).Encode(artworks)
		return
	}

	queryParams := query.NewArtworkQuery(r.URL.Query())
	artworks, err = h.store.FindMany(r.Context(), queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	err = h.cache.SetMany(r.Context(), queryString, artworks)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(artworks)
}
