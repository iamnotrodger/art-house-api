package exhibition

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
	router.HandleFunc("/api/exhibition", h.GetMany).Methods("GET")
	router.HandleFunc("/api/exhibition/{id}", h.Get).Methods("GET")
	router.HandleFunc("/api/exhibition/{id}/artwork", h.GetArtworks).Methods("GET")
	router.HandleFunc("/api/exhibition/{id}/artist", h.GetArtists).Methods("GET")
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	exhibitionID := params["id"]

	exhibition, err := h.cache.Get(r.Context(), exhibitionID)
	if err != nil {
		log.Println(err)
	} else if exhibition != nil {
		json.NewEncoder(w).Encode(exhibition)
		return
	}

	exhibition, err = h.store.Find(r.Context(), exhibitionID)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	err = h.cache.Set(r.Context(), exhibitionID, exhibition)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(exhibition)
}

func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryString := r.URL.RawQuery
	exhibitions, err := h.cache.GetMany(r.Context(), queryString)
	if err != nil {
		log.Println(err)
	} else if exhibitions != nil {
		json.NewEncoder(w).Encode(exhibitions)
		return
	}

	queryParams := query.NewExhibitionQuery(r.URL.Query())
	exhibitions, err = h.store.FindMany(r.Context(), queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	err = h.cache.SetMany(r.Context(), queryString, exhibitions)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(exhibitions)
}

func (h *Handler) GetArtworks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	exhibitionID := params["id"]

	queryString := r.URL.RawQuery
	artworks, err := h.cache.GetArtworks(r.Context(), exhibitionID, queryString)
	if err != nil {
		log.Println(err)
	} else if artworks != nil {
		json.NewEncoder(w).Encode(artworks)
		return
	}

	queryParams := query.NewArtworkQuery(r.URL.Query())
	artworks, err = h.store.FindArtworks(r.Context(), exhibitionID, queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	err = h.cache.SetArtworks(r.Context(), queryString, exhibitionID, artworks)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(artworks)
}

func (h *Handler) GetArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	exhibitionID := params["id"]

	queryString := r.URL.RawQuery
	artists, err := h.cache.GetArtists(r.Context(), exhibitionID, queryString)
	if err != nil {
		log.Println(err)
	} else if artists != nil {
		json.NewEncoder(w).Encode(artists)
		return
	}

	queryParams := query.NewArtistQuery(r.URL.Query())
	artists, err = h.store.FindArtists(r.Context(), exhibitionID, queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	err = h.cache.SetArtists(r.Context(), queryString, exhibitionID, artists)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(artists)
}
