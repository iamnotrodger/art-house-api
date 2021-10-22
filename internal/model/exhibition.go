package model

import (
	"context"

	"github.com/iamnotrodger/art-house-api/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Exhibition struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Images   []Image            `json:"images,omitempty" bson:"images,omitempty"`
	Artists  []Artist           `json:"artists,omitempty" bson:"artists,omitempty"`
	Artworks []Artwork          `json:"artworks,omitempty" bson:"artworks,omitempty"`
}

func (e *Exhibition) ConvertToBson() bson.D {
	var doc bson.D
	var artists []primitive.ObjectID
	var artworks []primitive.ObjectID

	for _, a := range e.Artists {
		artists = append(artists, a.ID)
	}

	for _, a := range e.Artworks {
		artworks = append(artworks, a.ID)
	}

	if !e.ID.IsZero() {
		doc = append(doc, bson.E{Key: "_id", Value: e.ID})
	}

	doc = append(doc,
		bson.E{Key: "name", Value: e.Name},
		bson.E{Key: "image", Value: e.Images},
		bson.E{Key: "artists", Value: artists},
		bson.E{Key: "artworks", Value: artworks},
	)

	return doc
}

func FindExhibition(db *mongo.Database, id primitive.ObjectID) (*Exhibition, error) {
	var exhibition Exhibition

	cursor, err := db.Collection("exhibitions").Find(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	cursor.Next(context.TODO())
	cursor.Decode(&exhibition)

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &exhibition, nil
}

func FindExhibitions(db *mongo.Database, filter bson.D, options ...*options.FindOptions) ([]Exhibition, error) {
	var exhibitions []Exhibition

	cursor, err := db.Collection("exhibitions").Find(context.TODO(), filter, options...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var exhibit Exhibition
		cursor.Decode(&exhibit)
		exhibitions = append(exhibitions, exhibit)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return exhibitions, nil
}

func FindExhibitionArtworks(db *mongo.Database, id primitive.ObjectID, options ...bson.D) ([]Artwork, error) {
	var exhibition Exhibition

	match := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
	unset := bson.D{{Key: "$unset", Value: "artists"}}
	lookupOptions := bson.D{
		{Key: "from", Value: "artworks"},
		{Key: "let", Value: bson.D{{
			Key: "artwork_ids", Value: "$artworks",
		}}},
		{Key: "as", Value: "artworks"},
	}
	lookupPipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{
			Key: "$expr", Value: bson.D{{
				Key: "$in", Value: bson.A{"$_id", "$$artwork_ids"},
			}},
		}},
		}},
		util.ArtworkLookup,
		util.ArtworkUnwind,
	}

	for _, option := range options {
		lookupPipeline = append(lookupPipeline, option)
	}

	lookupOptions = append(lookupOptions, bson.E{Key: "pipeline", Value: lookupPipeline})
	lookup := bson.D{
		{Key: "$lookup",
			Value: lookupOptions,
		}}

	pipeline := mongo.Pipeline{match, unset, lookup}

	cursor, err := db.Collection("exhibitions").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	cursor.Next(context.TODO())
	cursor.Decode(&exhibition)

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return exhibition.Artworks, nil
}

func FindExhibitionArtists(db *mongo.Database, id primitive.ObjectID, options ...bson.D) ([]Artist, error) {
	var exhibition Exhibition

	match := bson.D{{Key: "$match", Value: bson.M{"_id": id}}}
	lookupOptions := bson.D{
		{Key: "from", Value: "artists"},
		{Key: "let", Value: bson.D{{
			Key: "artist_ids", Value: "$artists",
		}}},
		{Key: "as", Value: "artists"},
	}
	lookupPipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{
			Key: "$expr", Value: bson.D{{
				Key: "$in", Value: bson.A{"$_id", "$$artist_ids"},
			}},
		}},
		}},
	}

	for _, option := range options {
		lookupPipeline = append(lookupPipeline, option)
	}

	lookupOptions = append(lookupOptions, bson.E{Key: "pipeline", Value: lookupPipeline})
	lookup := bson.D{
		{Key: "$lookup",
			Value: lookupOptions,
		}}

	pipeline := mongo.Pipeline{match, lookup}

	cursor, err := db.Collection("exhibitions").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if cursor.RemainingBatchLength() < 1 {
		return nil, mongo.ErrNoDocuments
	}

	cursor.Next(context.TODO())
	cursor.Decode(&exhibition)

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return exhibition.Artists, nil
}

func InsertExhibitions(db *mongo.Database, exhibitions []Exhibition) (*mongo.InsertManyResult, error) {
	var docs []interface{}

	for _, exhibit := range exhibitions {
		docs = append(docs, exhibit.ConvertToBson())
	}

	res, err := db.Collection("exhibitions").InsertMany(context.TODO(), docs)
	return res, err
}
