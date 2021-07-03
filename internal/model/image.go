package model

type Image struct {
	AspectRation string `json:"aspect_ratio,omitempty" bson:"aspect_ratio,omitempty"`
	Small        string `json:"small,omitempty" bson:"small,omitempty"`
	Large        string `json:"large,omitempty" bson:"large,omitempty"`
}
