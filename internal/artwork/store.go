package artwork

import (
	"context"

	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/iamnotrodger/art-house-api/internal/util"
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

func (s *Store) Find(id primitive.ObjectID) (*model.Artwork, error) {
	var artwork model.Artwork

	match := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
	limit := bson.D{{Key: "$limit", Value: 1}}

	pipeline := mongo.Pipeline{match, limit, util.ArtworkLookup, util.ArtworkUnwind}
	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	cursor.Next(context.TODO())
	cursor.Decode(&artwork)

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	model.SortImages(artwork.Images)
	return &artwork, nil
}

func (s *Store) FindMany(options ...bson.D) ([]model.Artwork, error) {
	var artworks = []model.Artwork{}

	unset := bson.D{{Key: "$unset", Value: "description"}}

	pipeline := mongo.Pipeline{}

	if index := util.FindLimitQuery(options); index > -1 {
		limit := options[index]
		options = append(options[:index], options[index+1:]...)
		pipeline = append(pipeline, options...)
		pipeline = append(pipeline, unset, util.ArtworkLookup, util.ArtworkUnwind)
		pipeline = append(pipeline, limit)
	} else {
		pipeline = append(pipeline, options...)
		pipeline = append(pipeline, unset, util.ArtworkLookup, util.ArtworkUnwind)
	}

	cursor, err := s.collection.Aggregate(context.TODO(), pipeline)
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

	for _, artwork := range artworks {
		model.SortImages(artwork.Images)
	}

	return artworks, nil
}

func (s *Store) InsertMany(artworks []model.Artwork) (*mongo.InsertManyResult, error) {
	var docs []interface{}

	for _, artwork := range artworks {
		model.SortImages(artwork.Images)
		docs = append(docs, artwork.ConvertToBson())
	}

	res, err := s.collection.InsertMany(context.TODO(), docs)
	return res, err
}
