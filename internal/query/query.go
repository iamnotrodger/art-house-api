package query

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryParams interface {
	GetFilter() bson.D
	GetFindOptions() *options.FindOptions
	GetPipeline() []bson.D
}
