package exhibition

import (
	"context"

	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		collection: db.Collection("exhibitions"),
	}
}

func (s *Store) Find(ctx context.Context, exhibitionID string) (*model.Exhibition, error) {
	id, err := primitive.ObjectIDFromHex(exhibitionID)
	if err != nil {
		return nil, primitive.ErrInvalidHex
	}

	cursor, err := s.collection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	var exhibition model.Exhibition
	cursor.Next(ctx)
	cursor.Decode(&exhibition)
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	model.SortImages(exhibition.Images)

	return &exhibition, nil
}

func (s *Store) FindMany(ctx context.Context, queryParam ...query.QueryParams) ([]*model.Exhibition, error) {
	var opts *options.FindOptions
	filter := bson.D{}

	if len(queryParam) > 0 {
		filter = queryParam[0].GetFilter()
		opts = queryParam[0].GetFindOptions()
	}

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	exhibitions := []*model.Exhibition{}
	for cursor.Next(ctx) {
		var exhibit model.Exhibition
		cursor.Decode(&exhibit)
		exhibitions = append(exhibitions, &exhibit)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	for _, exhibit := range exhibitions {
		model.SortImages(exhibit.Images)
	}

	return exhibitions, nil
}

func (s *Store) FindArtworks(ctx context.Context, exhibitionID string, queryParam ...query.QueryParams) ([]*model.Artwork, error) {
	id, err := primitive.ObjectIDFromHex(exhibitionID)
	if err != nil {
		return nil, primitive.ErrInvalidHex
	}

	match := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
	matchArtworks := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key: "$expr",
			Value: bson.D{{
				Key:   "$in",
				Value: bson.A{"$_id", "$$artwork_ids"},
			}},
		}},
	}}

	lookupPipeline := bson.A{matchArtworks}
	if len(queryParam) > 0 {
		for _, queryOpts := range queryParam[0].GetPipeline() {
			lookupPipeline = append(lookupPipeline, queryOpts)
		}
	}
	lookupPipeline = append(lookupPipeline, query.ArtworkLookupStage, query.ArtworkUnwindStage)

	lookup := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{Key: "from", Value: "artworks"},
			{Key: "let", Value: bson.D{{Key: "artwork_ids", Value: "$artwork_ids"}}},
			{Key: "pipeline", Value: lookupPipeline},
			{Key: "as", Value: "artworks"},
		},
	}}

	pipeline := mongo.Pipeline{match, lookup}
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	exhibition := model.Exhibition{}
	cursor.Next(ctx)
	cursor.Decode(&exhibition)
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	for _, artwork := range exhibition.Artworks {
		model.SortImages(artwork.Images)
		model.SortImages(artwork.Artist.Images)
	}

	return exhibition.Artworks, nil
}

func (s *Store) FindArtists(ctx context.Context, exhibitionID string, queryParam ...query.QueryParams) ([]*model.Artist, error) {
	id, err := primitive.ObjectIDFromHex(exhibitionID)
	if err != nil {
		return nil, primitive.ErrInvalidHex
	}

	match := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
	matchArtists := bson.D{{
		Key: "$match",
		Value: bson.D{{
			Key: "$expr",
			Value: bson.D{{
				Key:   "$in",
				Value: bson.A{"$_id", "$$artist_ids"},
			}},
		}},
	}}

	lookupPipeline := bson.A{matchArtists}
	if len(queryParam) > 0 {
		for _, queryOpts := range queryParam[0].GetPipeline() {
			lookupPipeline = append(lookupPipeline, queryOpts)
		}
	}

	lookup := bson.D{{
		Key: "$lookup",
		Value: bson.D{
			{Key: "from", Value: "artists"},
			{Key: "let", Value: bson.D{{Key: "artist_ids", Value: "$artist_ids"}}},
			{Key: "pipeline", Value: lookupPipeline},
			{Key: "as", Value: "artists"},
		},
	}}

	pipeline := mongo.Pipeline{match, lookup}
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	exhibition := model.Exhibition{}
	cursor.Next(ctx)
	cursor.Decode(&exhibition)
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	for _, artist := range exhibition.Artists {
		model.SortImages(artist.Images)
	}

	return exhibition.Artists, nil
}

func (s *Store) InsertMany(ctx context.Context, exhibitions []*model.Exhibition) error {
	var docs []interface{}

	for _, exhibit := range exhibitions {
		docs = append(docs, exhibit.ConvertToBson())
	}

	_, err := s.collection.InsertMany(ctx, docs)
	return err
}
