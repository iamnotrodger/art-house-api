package query

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ArtworkLookupStage = bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "artists"},
				{Key: "localField", Value: "artist_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "artist"},
			},
		}}

	ArtworkUnwindStage = bson.D{
		{Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$artist"},
				{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		}}
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
