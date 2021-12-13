package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	for _, artist := range e.Artists {
		artistID, err := primitive.ObjectIDFromHex(artist.ID)
		if err == nil {
			artists = append(artists, artistID)

		}
	}

	for _, artwork := range e.Artworks {
		artworks = append(artworks, artwork.ID)
	}

	if !e.ID.IsZero() {
		doc = append(doc, bson.E{Key: "_id", Value: e.ID})
	}

	doc = append(doc,
		bson.E{Key: "name", Value: e.Name},
		bson.E{Key: "images", Value: e.Images},
		bson.E{Key: "artists", Value: artists},
		bson.E{Key: "artworks", Value: artworks},
	)

	return doc
}
