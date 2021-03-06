package artwork

import (
	"context"
	"fmt"

	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		collection: db.Collection("artworks"),
	}
}

func (s *Store) Find(ctx context.Context, artworkID string) (*model.Artwork, error) {
	var artwork model.Artwork

	id, err := primitive.ObjectIDFromHex(artworkID)
	if err != nil {
		return nil, primitive.ErrInvalidHex
	}

	match := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
	limit := bson.D{{Key: "$limit", Value: 1}}

	pipeline := mongo.Pipeline{match, limit, query.ArtworkLookupStage, query.ArtworkUnwindStage}
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	cursor.Next(ctx)
	cursor.Decode(&artwork)

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	model.SortImages(artwork.Images)
	model.SortImages(artwork.Artist.Images)

	return &artwork, nil
}

func (s *Store) FindMany(ctx context.Context, queryParam ...query.QueryParams) ([]*model.Artwork, error) {
	pipeline := mongo.Pipeline{}

	if len(queryParam) > 0 {
		pipeline = queryParam[0].GetPipeline()
	}
	pipeline = append(pipeline, query.ArtworkLookupStage, query.ArtworkUnwindStage)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	artworks := []*model.Artwork{}
	err = cursor.All(ctx, &artworks)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal artworks: %w", err)
		return nil, err
	}

	for _, artwork := range artworks {
		model.SortImages(artwork.Images)
		model.SortImages(artwork.Artist.Images)
	}

	return artworks, nil
}

func (s *Store) InsertMany(ctx context.Context, artworks []*model.Artwork) error {
	var docs []interface{}

	for _, artwork := range artworks {
		model.SortImages(artwork.Images)
		docs = append(docs, artwork.ConvertToBson())
	}

	_, err := s.collection.InsertMany(ctx, docs)
	return err
}
