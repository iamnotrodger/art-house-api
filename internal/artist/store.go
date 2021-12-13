package artist

import (
	"context"
	"fmt"

	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		db:         db,
		collection: db.Collection("artists"),
	}
}

func (s *Store) Find(ctx context.Context, artistID string) (*model.Artist, error) {
	id, err := primitive.ObjectIDFromHex(artistID)
	if err != nil {
		return nil, primitive.ErrInvalidHex
	}

	singleRes := s.collection.FindOne(ctx, bson.M{"_id": id}, &options.FindOneOptions{})
	if err = singleRes.Err(); err != nil {
		return nil, err
	}

	artist := &model.Artist{}
	err = singleRes.Decode(artist)
	if err != nil {
		err = fmt.Errorf("error decoding artist: %w", err)
		return nil, err
	}

	model.SortImages(artist.Images)
	return artist, nil
}

func (s *Store) FindMany(ctx context.Context, filter bson.D, options ...*options.FindOptions) ([]*model.Artist, error) {
	cursor, err := s.collection.Find(ctx, filter, options...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	artists := []*model.Artist{}
	err = cursor.All(ctx, &artists)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal artists: %w", err)
		return nil, err
	}

	for _, artist := range artists {
		model.SortImages(artist.Images)
	}

	return artists, nil
}

func (s *Store) FindArtworks(ctx context.Context, artistID string, options ...bson.D) ([]*model.Artwork, error) {
	id, err := primitive.ObjectIDFromHex(artistID)
	if err != nil {
		return nil, primitive.ErrInvalidHex
	}

	match := bson.D{{Key: "$match", Value: bson.M{"artist": id}}}
	unset := bson.D{{Key: "$unset", Value: "description"}}

	pipeline := mongo.Pipeline{match}
	pipeline = append(pipeline, options...)
	pipeline = append(pipeline, unset, util.ArtworkLookup, util.ArtworkUnwind)

	cursor, err := s.db.Collection("artworks").Aggregate(ctx, pipeline)
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

	return artworks, nil
}

func (s *Store) InsertMany(ctx context.Context, artists []model.Artist) (*mongo.InsertManyResult, error) {
	var docs []interface{}

	for _, artist := range artists {
		model.SortImages(artist.Images)
		docs = append(docs, artist)
	}

	res, err := s.collection.InsertMany(ctx, docs)
	return res, err
}
