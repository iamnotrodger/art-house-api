package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Artwork struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Image       *Image             `json:"image,omitempty" bson:"title,omitempty"`
	Year        int                `json:"year,omitempty" bson:"year, omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Artist      *Artist            `json:"artist,omitempty" bson:"artist,omitempty"`
}
