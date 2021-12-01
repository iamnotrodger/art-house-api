package util

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	InvalidIDError = errors.New("Invalid ID")
)

func HandleError(w http.ResponseWriter, err error) {
	var statusCode int

	switch err {
	case mongo.ErrNilDocument:
		statusCode = http.StatusBadRequest
	case InvalidIDError:
		statusCode = http.StatusUnprocessableEntity
	default:
		statusCode = http.StatusInternalServerError
	}

	RespondWithError(w, statusCode, err.Error())
}

func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{ "message": "` + msg + `"}`))
}
