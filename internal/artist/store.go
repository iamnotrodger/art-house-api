package artist

import (
	"context"
	"fmt"

	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/query"
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

func (s *Store) FindMany(ctx context.Context, queryParam ...query.QueryParams) ([]*model.Artist, error) {
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

func (s *Store) FindArtworks(ctx context.Context, artistID string, queryParam ...query.QueryParams) ([]*model.Artwork, error) {
	id, err := primitive.ObjectIDFromHex(artistID)
	if err != nil {
		return nil, primitive.ErrInvalidHex
	}

	var opts *options.FindOptions
	filter := bson.D{{Key: "artist_id", Value: id}}
	if len(queryParam) > 0 {
		filter = append(filter, queryParam[0].GetFilter()...)
		opts = queryParam[0].GetFindOptions()
	} else {
		opts = options.Find()
	}
	opts.SetProjection(bson.M{"artist": 0})

	cursor, err := s.db.Collection("artworks").Find(ctx, filter, opts)
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
	}

	return artworks, nil
}

func (s *Store) InsertMany(ctx context.Context, artists []*model.Artist) error {
	var docs []interface{}

	for _, artist := range artists {
		model.SortImages(artist.Images)
		docs = append(docs, artist)
	}

	_, err := s.collection.InsertMany(ctx, docs)
	return err
}
