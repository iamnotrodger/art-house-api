package model

import (
	"sort"
)

type Image struct {
	Height *float64 `json:"height" bson:"height,omitempty"`
	Width  *float64 `json:"width" bson:"width,omitempty"`
	Url    string   `json:"url" bson:"url,omitempty"`
}

// SortImages sort the images by the size in decending order
func SortImages(images []*Image) {
	sort.Slice(images, func(i int, j int) bool {
		if images[i].Width == nil {
			return false
		} else if images[j].Width == nil {
			return true
		}

		return *images[i].Width < *images[j].Width
	})
}
