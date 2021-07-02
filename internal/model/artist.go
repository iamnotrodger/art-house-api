package model

import (
	"context"

	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Artist struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Image *Image             `json:"image,omitempty" bson:"image,omitempty"`
}

func GetArtists(db *mongo.Database) ([]Artist, error) {
	var artists = []Artist{}

	cursor, err := db.Collection("artists").Find(context.TODO(), bson.M{})
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

	return artists, nil
}

func GetArtistArtworks(db *mongo.Database, id primitive.ObjectID) ([]Artwork, error) {
	var artworks = []Artwork{}

	match := bson.D{{Key: "$match", Value: bson.M{"artist": id}}}
	unset := bson.D{{Key: "$unset", Value: "description"}}

	pipeline := mongo.Pipeline{match, unset, util.ArtworkLookup, util.ArtworkUnwind}
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
