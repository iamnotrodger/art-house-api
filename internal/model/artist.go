package model

import (
	"context"

	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Artist struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Images []Image            `json:"images,omitempty" bson:"images,omitempty"`
}

func FindArtist(db *mongo.Database, id primitive.ObjectID) (*Artist, error) {
	var artist Artist

	cursor, err := db.Collection("artists").Find(context.TODO(), bson.M{"_id": id})
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

	SortImages(artist.Images)
	return &artist, nil
}

func FindArtists(db *mongo.Database, filter bson.D, options ...*options.FindOptions) ([]Artist, error) {
	var artists = []Artist{}

	cursor, err := db.Collection("artists").Find(context.TODO(), filter, options...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var artist Artist
		cursor.Decode(&artist)
		artists = append(artists, artist)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	for _, artist := range artists {
		SortImages(artist.Images)
	}

	return artists, nil
}

func FindArtistArtworks(db *mongo.Database, id primitive.ObjectID, options ...bson.D) ([]Artwork, error) {
	var artworks = []Artwork{}

	match := bson.D{{Key: "$match", Value: bson.M{"artist": id}}}
	unset := bson.D{{Key: "$unset", Value: "description"}}

	pipeline := mongo.Pipeline{match}
	pipeline = append(pipeline, options...)
	pipeline = append(pipeline, unset, util.ArtworkLookup, util.ArtworkUnwind)

	cursor, err := db.Collection("artworks").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var artwork Artwork
		cursor.Decode(&artwork)
		artworks = append(artworks, artwork)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return artworks, nil
}

func InsertArtists(db *mongo.Database, artists []Artist) (*mongo.InsertManyResult, error) {
	var docs []interface{}

	for _, artist := range artists {
		SortImages(artist.Images)
		docs = append(docs, artist)
	}

	res, err := db.Collection("artists").InsertMany(context.TODO(), docs)
	return res, err
}
