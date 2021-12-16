package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Artist struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Images []*Image           `json:"images,omitempty" bson:"images,omitempty"`
}
