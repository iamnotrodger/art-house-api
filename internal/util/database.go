package util

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/iamnotrodger/art-house-api/cmd/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Global.MongoURI))
	if err != nil {
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Global.RedisAddr,
		Password: config.Global.RedisPassword,
		DB:       config.Global.RedisDb,
	})
}
