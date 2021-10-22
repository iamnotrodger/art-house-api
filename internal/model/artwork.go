package model

import (
	"context"

	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Artwork struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Images      []Image            `json:"images,omitempty" bson:"images,omitempty"`
	Year        int                `json:"year,omitempty" bson:"year, omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Artist      *Artist            `json:"artist,omitempty" bson:"artist,omitempty"`
}

func (a *Artwork) ConvertToBson() bson.D {
	var doc bson.D

	if !a.ID.IsZero() {
		doc = append(doc, bson.E{Key: "_id", Value: a.ID})
	}

	doc = append(doc,
		bson.E{Key: "title", Value: a.Title},
		bson.E{Key: "images", Value: a.Images},
		bson.E{Key: "year", Value: a.Year},
		bson.E{Key: "description", Value: a.Description},
		bson.E{Key: "artist", Value: a.Artist.ID},
	)

	return doc
}

func FindArtwork(db *mongo.Database, id primitive.ObjectID) (*Artwork, error) {
	var artwork Artwork

	match := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
	limit := bson.D{{Key: "$limit", Value: 1}}

	pipeline := mongo.Pipeline{match, limit, util.ArtworkLookup, util.ArtworkUnwind}
	cursor, err := db.Collection("artworks").Aggregate(context.TODO(), pipeline)
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

	SortImages(artwork.Images)
	return &artwork, nil
}

func FindArtworks(db *mongo.Database, options ...bson.D) ([]Artwork, error) {
	var artworks = []Artwork{}

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

	for _, artwork := range artworks {
		SortImages(artwork.Images)
	}

	return artworks, nil
}

func InsertArtworks(db *mongo.Database, artworks []Artwork) (*mongo.InsertManyResult, error) {
	var docs []interface{}

	for _, artwork := range artworks {
		SortImages(artwork.Images)
		docs = append(docs, artwork.ConvertToBson())
	}

	res, err := db.Collection("artworks").InsertMany(context.TODO(), docs)
	return res, err
}
