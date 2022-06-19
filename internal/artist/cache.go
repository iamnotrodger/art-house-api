package artist

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
		namespace:  "artist",
	}
}

func (c *Cache) Get(ctx context.Context, artistID string) (*model.Artist, error) {
	key := c.getKeyByID(artistID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var artist model.Artist
	err = json.Unmarshal([]byte(val), &artist)
	if err != nil {
		return nil, err
	}

	return &artist, nil
}

func (c *Cache) GetMany(ctx context.Context, queryString string) ([]*model.Artist, error) {
	key := c.getKeyByQuery(queryString)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var artists []*model.Artist
	err = json.Unmarshal([]byte(val), &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (c *Cache) GetArtworks(ctx context.Context, artistID string, queryString string) ([]*model.Artwork, error) {
	key := c.getKeyByArtworks(artistID, queryString)
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

func (c *Cache) Set(ctx context.Context, artistID string, artist *model.Artist) error {
	artistJson, err := json.Marshal(artist)
	if err != nil {
		return err
	}

	key := c.getKeyByID(artistID)
	err = c.client.Set(ctx, key, artistJson, c.expiration).Err()
	return err
}

func (c *Cache) SetMany(ctx context.Context, queryString string, artists []*model.Artist) error {
	artistsJson, err := json.Marshal(artists)
	if err != nil {
		return err
	}

	key := c.getKeyByQuery(queryString)
	err = c.client.Set(ctx, key, artistsJson, c.expiration).Err()
	return err
}

func (c *Cache) SetArtworks(ctx context.Context, artistID string, queryString string, artworks []*model.Artwork) error {
	artworksJson, err := json.Marshal(artworks)
	if err != nil {
		return err
	}

	key := c.getKeyByArtworks(artistID, queryString)
	err = c.client.Set(ctx, key, artworksJson, c.expiration).Err()
	return err
}

func (c *Cache) getKeyByID(artistID string) string {
	return fmt.Sprintf("%s:%s", c.namespace, artistID)
}

func (c *Cache) getKeyByQuery(queryString string) string {
	return fmt.Sprintf("%s?%s", c.namespace, queryString)
}

func (c *Cache) getKeyByArtworks(artistID string, queryString string) string {
	return fmt.Sprintf("%s:%s:%s?%s", c.namespace, artistID, "artwork", queryString)
}
