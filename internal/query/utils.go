package query

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func parseSort(sortString string) (string, int) {
	pair := strings.Split(sortString, ":")
	if len(pair) != 2 {
		return "", 0
	}

	key := pair[0]
	order := pair[1]
	value := 0

	if order == "asc" {
		value = 1
	} else if order == "desc" {
		value = -1
	} else {
		return "", 0
	}

	return key, value
}

func getSortAsBson(sortMap map[string]int) bson.D {
	sort := bson.D{}
	for key, value := range sortMap {
		sort = append(sort, bson.E{Key: key, Value: value})
	}
	return sort
}
