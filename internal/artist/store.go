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

func (s *Store) Find(id primitive.ObjectID) (*model.Artist, error) {
	var artist model.Artist

	cursor, err := s.collection.Find(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	cursor.Next(context.TODO())
	cursor.Decode(&artist)

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	model.SortImages(artist.Images)
	return &artist, nil
}

func (s *Store) FindMany(filter bson.D, options ...*options.FindOptions) ([]model.Artist, error) {
	var artists = []model.Artist{}

	cursor, err := s.collection.Find(context.TODO(), filter, options...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
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

func (s *Store) FindArtworks(id primitive.ObjectID, options ...bson.D) ([]model.Artwork, error) {
	var artworks = []model.Artwork{}

	match := bson.D{{Key: "$match", Value: bson.M{"artist": id}}}
	unset := bson.D{{Key: "$unset", Value: "description"}}

	pipeline := mongo.Pipeline{match}
	pipeline = append(pipeline, options...)
	pipeline = append(pipeline, unset, util.ArtworkLookup, util.ArtworkUnwind)

	cursor, err := s.db.Collection("artworks").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var artwork model.Artwork
		cursor.Decode(&artwork)
		artworks = append(artworks, artwork)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return artworks, nil
}

func (s *Store) InsertMany(artists []model.Artist) (*mongo.InsertManyResult, error) {
	var docs []interface{}

	for _, artist := range artists {
		model.SortImages(artist.Images)
		docs = append(docs, artist)
	}

	res, err := s.collection.InsertMany(context.TODO(), docs)
	return res, err
}
