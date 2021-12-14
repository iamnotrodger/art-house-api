package model

type Artist struct {
	ID     string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string   `json:"name,omitempty" bson:"name,omitempty"`
	Images []*Image `json:"images,omitempty" bson:"images,omitempty"`
}
