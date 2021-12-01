package artist

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
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

	options := util.QueryBuilder(r.URL.Query())
	artists, err := h.store.FindMany(r.Context(), bson.D{}, options)
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

	queryParams := r.URL.Query()
	delete(queryParams, "search")

	options := util.QueryBuilderPipeline(queryParams)
	artworks, err := h.store.FindArtworks(r.Context(), artistID, options...)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(artworks)
}
