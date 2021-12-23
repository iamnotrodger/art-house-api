package query

import (
	"strconv"

	"github.com/iamnotrodger/art-house-api/cmd/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExhibitionQueryParams struct {
	limit int64
	skip  int64
	sort  map[string]int
}

func NewExhibitionQuery(parameters map[string][]string) *ExhibitionQueryParams {
	query := &ExhibitionQueryParams{}

	if limit, ok := parameters["limit"]; ok {
		query.setLimitFromString(limit[0])
	} else {
		query.limit = config.Global.ExhibitionLimit
	}
	if skip, ok := parameters["skip"]; ok {
		query.setSkipFromString(skip[0])
	}
	if sort, ok := parameters["sort"]; ok {
		query.SetSort(sort)
	}

	return query
}

func (q *ExhibitionQueryParams) GetFilter() bson.D {
	return bson.D{}
}

func (q *ExhibitionQueryParams) GetFindOptions() *options.FindOptions {
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

func (q *ExhibitionQueryParams) GetPipeline() []bson.D {
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

func (q *ExhibitionQueryParams) SetLimit(limit int64) {
	if limit < config.Global.ExhibitionLimitMin {
		q.limit = config.Global.ExhibitionLimit
	} else if limit > config.Global.ExhibitionLimitMax {
		q.limit = config.Global.ExhibitionLimitMax
	} else {
		q.limit = limit
	}
}
func (q *ExhibitionQueryParams) SetSkip(skip int64) {
	if skip > 0 {
		q.skip = skip
	}
}

func (q *ExhibitionQueryParams) SetSort(sortArray []string) {
	q.sort = map[string]int{}
	for _, sortString := range sortArray {
		key, value := parseSort(sortString)
		if key != "" {
			q.sort[key] = value
		}
	}
}

func (q *ExhibitionQueryParams) setLimitFromString(limitString string) {
	limit, err := strconv.ParseInt(limitString, 0, 64)
	if err != nil {
		q.limit = config.Global.ExhibitionLimit
	} else {
		q.SetLimit(limit)
	}
}

func (q *ExhibitionQueryParams) setSkipFromString(skipString string) {
	skip, err := strconv.ParseInt(skipString, 0, 64)
	if err == nil {
		q.SetSkip(skip)
	}
}

func (q *ExhibitionQueryParams) isSortValid() bool {
	return q.sort != nil && len(q.sort) > 0
}

func (q *ExhibitionQueryParams) isSkipValid() bool {
	return q.skip > 0
}
