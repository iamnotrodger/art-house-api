package query

import (
	"strconv"

	"github.com/iamnotrodger/art-house-api/cmd/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArtworkQueryParams struct {
	limit  int64
	skip   int64
	sort   map[string]int
	search string
}

func NewArtworkQuery(parameters map[string][]string) *ArtworkQueryParams {
	query := &ArtworkQueryParams{}

	if limit, ok := parameters["limit"]; ok {
		query.setLimitFromString(limit[0])
	} else {
		query.limit = config.Global.ArtworkLimit
	}
	if skip, ok := parameters["skip"]; ok {
		query.setSkipFromString(skip[0])
	}
	if sort, ok := parameters["sort"]; ok {
		query.SetSort(sort)
	}
	if search, ok := parameters["search"]; ok {
		query.SetSearch(search[0])
	}

	return query
}

func (q *ArtworkQueryParams) GetFilter() bson.D {
	filter := bson.D{}
	if q.isSearchValid() {
		search := bson.D{{Key: "$search", Value: q.search}}
		text := bson.D{{Key: "$text", Value: search}}
		filter = append(filter, text...)
	}
	return filter
}

func (q *ArtworkQueryParams) GetFindOptions() *options.FindOptions {
	options := options.Find()
	if q.isSortValid() {
		sort := getSortAsBson(q.sort)
		options.SetSort(sort)
	}
	if q.isSkipValid() {
		options.SetSkip(q.skip)
	}
	options.SetLimit(q.limit)
	return options
}

func (q *ArtworkQueryParams) GetPipeline() []bson.D {
	pipeline := []bson.D{}

	if q.isSearchValid() {
		search := bson.D{{Key: "$search", Value: q.search}}
		text := bson.D{{Key: "$text", Value: search}}
		match := bson.D{{Key: "$match", Value: text}}
		pipeline = append(pipeline, match)
	}
	if q.isSortValid() {
		sort := bson.D{{Key: "$sort", Value: q.sort}}
		pipeline = append(pipeline, sort)
	}
	if q.isSkipValid() {
		skip := bson.D{{Key: "$skip", Value: q.skip}}
		pipeline = append(pipeline, skip)
	}
	limit := bson.D{{Key: "$limit", Value: q.limit}}
	pipeline = append(pipeline, limit)

	return pipeline
}

func (q *ArtworkQueryParams) SetLimit(limit int64) {
	if limit < config.Global.ArtworkLimitMin {
		q.limit = config.Global.ArtworkLimit
	} else if limit > config.Global.ArtworkLimitMax {
		q.limit = config.Global.ArtworkLimitMax
	} else {
		q.limit = limit
	}
}
func (q *ArtworkQueryParams) SetSkip(skip int64) {
	if skip > 0 {
		q.skip = skip
	}
}

func (q *ArtworkQueryParams) SetSort(sortArray []string) {
	q.sort = map[string]int{}
	for _, sortString := range sortArray {
		key, value := parseSort(sortString)
		if key != "" {
			q.sort[key] = value
		}
	}
}

func (q *ArtworkQueryParams) SetSearch(search string) {
	if search != "" {
		q.search = search
	}
}

func (q *ArtworkQueryParams) setLimitFromString(limitString string) {
	limit, err := strconv.ParseInt(limitString, 0, 64)
	if err != nil {
		q.limit = config.Global.ArtworkLimit
	} else {
		q.SetLimit(limit)
	}
}

func (q *ArtworkQueryParams) setSkipFromString(skipString string) {
	skip, err := strconv.ParseInt(skipString, 0, 64)
	if err == nil {
		q.SetSkip(skip)
	}
}

func (q *ArtworkQueryParams) isSortValid() bool {
	return q.sort != nil && len(q.sort) > 0
}

func (q *ArtworkQueryParams) isSkipValid() bool {
	return q.skip > 0
}

func (q *ArtworkQueryParams) isSearchValid() bool {
	return q.search != ""
}
