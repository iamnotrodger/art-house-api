package util

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ArtworkLookup = bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "artists"},
				{Key: "localField", Value: "artist"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "artist"},
			},
		}}

	ArtworkUnwind = bson.D{
		{Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$artist"},
				{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		}}
)

type pairInt struct {
	Key   string
	Value int
}

func QueryBuilder(parameters url.Values) *options.FindOptions {
	options := options.Find()

	if sortArray, ok := parameters["sort"]; ok {
		var sort bson.D

		for _, sortString := range sortArray {
			pair, err := parseSort(sortString)
			if err == nil {
				sort = append(sort, bson.E{Key: pair.Key, Value: pair.Value})
			}
		}

		if sort != nil {
			options.SetSort(sort)
		}
	}

	if skipString, ok := parameters["skip"]; ok {
		skip, err := strconv.ParseInt(skipString[0], 0, 64)
		if err == nil && skip > 0 {
			options.SetSkip(skip)
		}
	}

	if limitString, ok := parameters["limit"]; ok {
		limit, err := strconv.ParseInt(limitString[0], 0, 64)
		if err == nil && limit > 0 {
			options.SetLimit(limit)
		}
	}

	return options
}

//TODO: update this function so that it doesn't affect the url.Values
//delete doesn't seem like a good idea
func QueryBuilderPipeline(parameters url.Values) []bson.D {
	var options []bson.D

	if sortArray, ok := parameters["sort"]; ok {
		var sort bson.D

		for _, sortString := range sortArray {
			pair, err := parseSort(sortString)
			if err == nil {
				sort = append(sort, bson.E{Key: pair.Key, Value: pair.Value})
			}
		}

		if sort != nil {
			options = append(options, bson.D{{Key: "$sort", Value: sort}})
		}
		delete(parameters, "sort")
	}

	if skipString, ok := parameters["skip"]; ok {
		skip, err := strconv.ParseInt(skipString[0], 0, 64)
		if err == nil && skip > 0 {
			options = append(options, bson.D{{Key: "$skip", Value: skip}})
		}
		delete(parameters, "skip")
	}

	if limitString, ok := parameters["limit"]; ok {
		limit, err := strconv.ParseInt(limitString[0], 0, 64)
		if err == nil && limit > 0 {
			options = append(options, bson.D{{Key: "$limit", Value: limit}})
		}
		delete(parameters, "limit")
	}

	if len(parameters) > 0 {
		match := bson.M{}

		if search, ok := parameters["search"]; ok {
			match["$text"] = bson.D{{Key: "$search", Value: search[0]}}
			delete(parameters, "search")
		}

		for key, value := range parameters {
			match[key] = value[0]
		}
		options = append([]bson.D{{{Key: "$match", Value: match}}}, options...)
	}

	return options
}

// FindOptionIndex returns -1 if index is not found
func FindOptionIndex(key string, options []bson.D) int {
	for i := len(options) - 1; i >= 0; i-- {
		optionMap := options[i].Map()
		if _, ok := optionMap[key]; ok {
			return i
		}
	}
	return -1
}

func FindLimitQuery(options []bson.D) int {
	for i := len(options) - 1; i >= 0; i-- {
		optionMap := options[i].Map()
		if _, ok := optionMap["$limit"]; ok {
			return i
		}
	}

	return -1
}

func parseSort(sort string) (*pairInt, error) {
	sortPair := &pairInt{
		Value: 1,
	}

	pair := strings.Split(sort, ":")
	if len(pair) != 2 {
		return nil, errors.New("invalid sort pair")
	}

	sortPair.Key = pair[0]
	if pair[1] == "desc" {
		sortPair.Value = -1
	}

	return sortPair, nil
}
