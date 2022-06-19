package artist

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/internal/query"
	"github.com/iamnotrodger/art-house-api/internal/util"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{
		store: store,
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

	artist, err := h.store.Find(r.Context(), artistID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(artist)
}

func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParams := query.NewArtistQuery(r.URL.Query())
	artists, err := h.store.FindMany(r.Context(), queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(artists)
}

func (h *Handler) GetArtworks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	artistID := params["id"]

	queryParams := query.NewArtworkQuery(r.URL.Query())
	artworks, err := h.store.FindArtworks(r.Context(), artistID, queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	artist, err := h.store.Find(r.Context(), artistID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	for _, artwork := range artworks {
		artwork.Artist = artist
	}

	json.NewEncoder(w).Encode(artworks)
}
