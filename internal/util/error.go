package util

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func HandleError(w http.ResponseWriter, err error) {
	var statusCode int

	switch err {
	case mongo.ErrNilDocument:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	RespondWithError(w, statusCode, err.Error())
}

func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{ "message": "` + msg + `"}`))
}
