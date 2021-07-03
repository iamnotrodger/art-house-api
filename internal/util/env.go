package util

import (
	"errors"
	"os"
)

func GetPort() string {
	var port string

	if portEnv, ok := os.LookupEnv("PORT"); ok {
		port = ":" + portEnv
	} else {
		port = ":8080"
	}

	return port
}

func GetDatabaseURI() (string, error) {
	if url, ok := os.LookupEnv("DATABASE_URI"); ok {
		return url, nil
	} else {
		return "", errors.New("missing DATABASE_URI env")
	}

}

func GetDatabaseName() (string, error) {
	if name, ok := os.LookupEnv("DATABASE"); ok {
		return name, nil
	} else {
		return "", errors.New("missing DATABASE env")
	}
}
