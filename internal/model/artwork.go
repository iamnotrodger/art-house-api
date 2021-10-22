package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
