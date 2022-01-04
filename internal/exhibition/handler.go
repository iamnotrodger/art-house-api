package exhibition

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/internal/query"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	store *Store
}

func NewHandler(db *mongo.Database) *Handler {
	return &Handler{
		store: NewStore(db),
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

	exhibition, err := h.store.Find(r.Context(), exhibitionID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(exhibition)
}

func (h *Handler) GetMany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParams := query.NewExhibitionQuery(r.URL.Query())
	exhibitions, err := h.store.FindMany(r.Context(), queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(exhibitions)
}

func (h *Handler) GetArtworks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	exhibitionID := params["id"]

	queryParams := query.NewArtworkQuery(r.URL.Query())
	artworks, err := h.store.FindArtworks(r.Context(), exhibitionID, queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(artworks)
}

func (h *Handler) GetArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	exhibitionID := params["id"]

	queryParams := query.NewArtistQuery(r.URL.Query())
	artists, err := h.store.FindArtists(r.Context(), exhibitionID, queryParams)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(artists)
}
