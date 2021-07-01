package handler

import (
	"fmt"
	"net/http"
)

func Health(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "OKAY")
}
