package artwork

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/iamnotrodger/art-house-api/internal/model"
)

type Cache struct {
	client     *redis.Client
	expiration time.Duration
	namespace  string
}

func NewCache(client *redis.Client, expiration time.Duration) *Cache {
	return &Cache{
		client:     client,
		expiration: expiration,
		namespace:  "artwork",
	}
}

func (c *Cache) Get(ctx context.Context, artworkID string) (*model.Artwork, error) {
	key := c.getKeyByID(artworkID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var artwork model.Artwork
	err = json.Unmarshal([]byte(val), &artwork)
	if err != nil {
		return nil, err
	}

	return &artwork, nil
}

func (c *Cache) GetMany(ctx context.Context, queryString string) ([]*model.Artwork, error) {
	key := c.getKeyByQuery(queryString)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var artworks []*model.Artwork
	err = json.Unmarshal([]byte(val), &artworks)
	if err != nil {
		return nil, err
	}

	return artworks, nil
}

func (c *Cache) Set(ctx context.Context, artworkID string, artwork *model.Artwork) error {
	artworkJson, err := json.Marshal(artwork)
	if err != nil {
		return err
	}

	key := c.getKeyByID(artworkID)
	err = c.client.Set(ctx, key, artworkJson, c.expiration).Err()
	return err
}

func (c *Cache) SetMany(ctx context.Context, queryString string, artworks []*model.Artwork) error {
	artworksJson, err := json.Marshal(artworks)
	if err != nil {
		return err
	}

	key := c.getKeyByQuery(queryString)
	err = c.client.Set(ctx, key, artworksJson, c.expiration).Err()
	return err
}

func (c *Cache) getKeyByID(artworkID string) string {
	return fmt.Sprintf("%s:%s", c.namespace, artworkID)
}

func (c *Cache) getKeyByQuery(queryString string) string {
	return fmt.Sprintf("%s?%s", c.namespace, queryString)
}
