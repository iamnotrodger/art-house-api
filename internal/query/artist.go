package query

import (
	"strconv"

	"github.com/iamnotrodger/art-house-api/cmd/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArtistQueryParams struct {
	limit int64
	skip  int64
	sort  map[string]int
}

func NewArtistQuery(parameters map[string][]string) *ArtistQueryParams {
	query := &ArtistQueryParams{}

	if limit, ok := parameters["limit"]; ok {
		query.setLimitFromString(limit[0])
	} else {
		query.limit = config.Global.ArtistLimit
	}
	if skip, ok := parameters["skip"]; ok {
		query.setSkipFromString(skip[0])
	}
	if sort, ok := parameters["sort"]; ok {
		query.SetSort(sort)
	}

	return query
}

func (q *ArtistQueryParams) GetFilter() bson.D {
	return bson.D{}
}

func (q *ArtistQueryParams) GetFindOptions() *options.FindOptions {
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

func (q *ArtistQueryParams) GetPipeline() []bson.D {
	pipeline := []bson.D{}

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

func (q *ArtistQueryParams) SetLimit(limit int64) {
	if limit < config.Global.ArtistLimitMin {
		q.limit = config.Global.ArtistLimit
	} else if limit > config.Global.ArtistLimitMax {
		q.limit = config.Global.ArtworkLimitMax
	} else {
		q.limit = limit
	}
}
func (q *ArtistQueryParams) SetSkip(skip int64) {
	if skip > 0 {
		q.skip = skip
	}
}

func (q *ArtistQueryParams) SetSort(sortArray []string) {
	q.sort = map[string]int{}
	for _, sortString := range sortArray {
		key, value := parseSort(sortString)
		if key != "" {
			q.sort[key] = value
		}
	}
}

func (q *ArtistQueryParams) setLimitFromString(limitString string) {
	limit, err := strconv.ParseInt(limitString, 0, 64)
	if err != nil {
		q.limit = config.Global.ArtworkLimit
	} else {
		q.SetLimit(limit)
	}
}

func (q *ArtistQueryParams) setSkipFromString(skipString string) {
	skip, err := strconv.ParseInt(skipString, 0, 64)
	if err == nil {
		q.SetSkip(skip)
	}
}

func (q *ArtistQueryParams) isSortValid() bool {
	return q.sort != nil && len(q.sort) > 0
}

func (q *ArtistQueryParams) isSkipValid() bool {
	return q.skip > 0
}
