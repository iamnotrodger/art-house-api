package model

type Image struct {
	Small string `json:"small,omitempty" bson:"small,omitempty"`
	Large string `json:"large,omitempty" bson:"large,omitempty"`
}
