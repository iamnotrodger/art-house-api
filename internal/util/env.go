package util

import "os"

func GetPort() string {
	var port string

	if portEnv, ok := os.LookupEnv("PORT"); ok {
		port = ":" + portEnv
	} else {
		port = ":8080"
	}

	return port
}
