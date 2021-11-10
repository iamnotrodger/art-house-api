package artist

import (
	"context"

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

func (s *Store) Find(ctx context.Context, id primitive.ObjectID) (*model.Artist, error) {
	var artist model.Artist

	cursor, err := s.collection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	cursor.Next(ctx)
	cursor.Decode(&artist)

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	model.SortImages(artist.Images)
	return &artist, nil
}

func (s *Store) FindMany(ctx context.Context, filter bson.D, options ...*options.FindOptions) ([]model.Artist, error) {
	var artists = []model.Artist{}

	cursor, err := s.collection.Find(ctx, filter, options...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var artist model.Artist
		cursor.Decode(&artist)
		artists = append(artists, artist)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	for _, artist := range artists {
		model.SortImages(artist.Images)
	}

	return artists, nil
}

func (s *Store) FindArtworks(ctx context.Context, id primitive.ObjectID, options ...bson.D) ([]model.Artwork, error) {
	var artworks = []model.Artwork{}

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

	for cursor.Next(ctx) {
		var artwork model.Artwork
		cursor.Decode(&artwork)
		artworks = append(artworks, artwork)
	}

	if err = cursor.Err(); err != nil {
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
