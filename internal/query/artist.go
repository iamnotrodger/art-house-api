package query

import (
	"strconv"

	"github.com/iamnotrodger/art-house-api/cmd/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArtistQueryParams struct {
	Limit int64
	Skip  int64
	Sort  map[string]int
}

func NewArtistQuery(parameters map[string][]string) *ArtistQueryParams {
	query := &ArtistQueryParams{}

	if limit, ok := parameters["limit"]; ok {
		query.setLimit(limit[0])
	} else {
		query.Limit = config.Global.ArtistLimit
	}
	if skip, ok := parameters["skip"]; ok {
		query.setSkip(skip[0])
	}
	if sort, ok := parameters["sort"]; ok {
		query.setSort(sort)
	}

	return query
}

func (q *ArtistQueryParams) GetFilter() bson.D {
	return bson.D{}
}

func (q *ArtistQueryParams) GetFindOptions() *options.FindOptions {
	options := options.Find()
	if q.isSortValid() {
		sort := q.getSortBson()
		options.SetSort(sort)
	}
	if q.isSkipValid() {
		options.SetSkip(q.Skip)
	}
	options.SetLimit(q.Limit)
	return options
}

func (q *ArtistQueryParams) GetPipeline() []bson.D {
	return nil
}

func (q *ArtistQueryParams) setLimit(limitString string) {
	limit, err := strconv.ParseInt(limitString, 0, 64)
	if err != nil || limit < config.Global.ArtworkLimitMin || limit > config.Global.ArtworkLimitMax {
		q.Limit = config.Global.ArtistLimit
	} else {
		q.Limit = limit
	}
}

func (q *ArtistQueryParams) setSkip(skipString string) {
	skip, err := strconv.ParseInt(skipString, 0, 64)
	if err == nil && skip > 0 {
		q.Skip = skip
	}
}

func (q *ArtistQueryParams) setSort(sortArray []string) {
	q.Sort = map[string]int{}
	for _, sortString := range sortArray {
		key, value := parseSort(sortString)
		if key != "" {
			q.Sort[key] = value
		}
	}

}

func (q *ArtistQueryParams) isSortValid() bool {
	return len(q.Sort) > 0
}

func (q *ArtistQueryParams) isSkipValid() bool {
	return q.Skip > 0
}

func (q *ArtistQueryParams) getSortBson() bson.D {
	sort := bson.D{}
	for key, value := range q.Sort {
		sort = append(sort, bson.E{Key: key, Value: value})
	}
	return sort
}
