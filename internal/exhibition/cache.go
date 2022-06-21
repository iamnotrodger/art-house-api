package exhibition

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
		namespace:  "exhibition",
	}
}

func (c *Cache) Get(ctx context.Context, exhibitionID string) (*model.Exhibition, error) {
	key := c.getKeyByID(exhibitionID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var exhibition model.Exhibition
	err = json.Unmarshal([]byte(val), &exhibition)
	if err != nil {
		return nil, err
	}

	return &exhibition, nil
}

func (c *Cache) GetMany(ctx context.Context, queryString string) ([]*model.Exhibition, error) {
	key := c.getKeyByQuery(queryString)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var exhibitions []*model.Exhibition
	err = json.Unmarshal([]byte(val), &exhibitions)
	if err != nil {
		return nil, err
	}

	return exhibitions, nil
}

func (c *Cache) GetArtworks(ctx context.Context, exhibitionID string, queryString string) ([]*model.Artwork, error) {
	key := c.getKeyByArtworks(exhibitionID, queryString)
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

func (c *Cache) GetArtists(ctx context.Context, exhibitionID string, queryString string) ([]*model.Artist, error) {
	key := c.getKeyByArtists(exhibitionID, queryString)
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

func (c *Cache) Set(ctx context.Context, exhibitionID string, exhibition *model.Exhibition) error {
	exhibitionJson, err := json.Marshal(exhibition)
	if err != nil {
		return err
	}

	key := c.getKeyByID(exhibitionID)
	err = c.client.Set(ctx, key, exhibitionJson, c.expiration).Err()
	return err
}

func (c *Cache) SetMany(ctx context.Context, queryString string, exhibitions []*model.Exhibition) error {
	exhibitionsJson, err := json.Marshal(exhibitions)
	if err != nil {
		return err
	}

	key := c.getKeyByQuery(queryString)
	err = c.client.Set(ctx, key, exhibitionsJson, c.expiration).Err()
	return err
}

func (c *Cache) SetArtworks(ctx context.Context, exhibitionID string, queryString string, artworks []*model.Artwork) error {
	artworksJson, err := json.Marshal(artworks)
	if err != nil {
		return err
	}

	key := c.getKeyByArtworks(exhibitionID, queryString)
	err = c.client.Set(ctx, key, artworksJson, c.expiration).Err()
	return err
}

func (c *Cache) SetArtists(ctx context.Context, exhibitionID string, queryString string, artworks []*model.Artist) error {
	artistJson, err := json.Marshal(artworks)
	if err != nil {
		return err
	}

	key := c.getKeyByArtists(exhibitionID, queryString)
	err = c.client.Set(ctx, key, artistJson, c.expiration).Err()
	return err
}

func (c *Cache) getKeyByID(exhibitionID string) string {
	return fmt.Sprintf("%s:%s", c.namespace, exhibitionID)
}

func (c *Cache) getKeyByQuery(queryString string) string {
	return fmt.Sprintf("%s?%s", c.namespace, queryString)
}

func (c *Cache) getKeyByArtworks(exhibitionID string, queryString string) string {
	return fmt.Sprintf("%s:%s:%s?%s", c.namespace, exhibitionID, "artwork", queryString)
}

func (c *Cache) getKeyByArtists(exhibitionID string, queryString string) string {
	return fmt.Sprintf("%s:%s:%s?%s", c.namespace, exhibitionID, "artist", queryString)
}
