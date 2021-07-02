package util

import "go.mongodb.org/mongo-driver/bson"

var ArtworkLookup = bson.D{
	{Key: "$lookup",
		Value: bson.D{
			{Key: "from", Value: "artists"},
			{Key: "localField", Value: "artist"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "artist"},
		},
	}}

var ArtworkUnwind = bson.D{
	{Key: "$unwind",
		Value: bson.D{
			{Key: "path", Value: "$artist"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		},
	}}
