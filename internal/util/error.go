package util

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleError(w http.ResponseWriter, err error) {
	var statusCode int
	message := err.Error()

	switch err {
	case mongo.ErrNilDocument:
		statusCode = http.StatusBadRequest
	case mongo.ErrNoDocuments:
		statusCode = http.StatusNotFound
	case primitive.ErrInvalidHex:
		statusCode = http.StatusUnprocessableEntity
		message = "Invalid ID"
	default:
		statusCode = http.StatusInternalServerError
	}

	RespondWithError(w, statusCode, message)
}

func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}
