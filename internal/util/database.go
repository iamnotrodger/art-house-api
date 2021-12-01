package util

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(ctx context.Context, databaseURI string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURI))
	if err != nil {
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
