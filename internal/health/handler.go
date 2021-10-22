package health

import (
	"fmt"
	"net/http"
)

func GetHealth(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "OKAY")
}
