package artwork

import (
	"context"
	"errors"
	"testing"

	"github.com/iamnotrodger/art-house-api/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	MongoFailResponse    = bson.D{{Key: "ok", Value: 0}}
	ErrMongoCommandError = mongo.CommandError{Message: "command failed"}
	ErrMongoNoResponses  = mongo.CommandError{Message: "no responses remaining", Labels: []string{"NetworkError"}, Wrapped: errors.New("no responses remaining")}

	artworkID          = "60e0850266d6c13d7b599b69"
	artistID           = "60e0850266d6c13d7b599b6a"
	artworkObjectID, _ = primitive.ObjectIDFromHex(artworkID)
	artistObjectID, _  = primitive.ObjectIDFromHex(artistID)
	imageSizeOne       = 1.0
	imageSizeTwo       = 2.0

	artist = &model.Artist{
		ID:     artistObjectID,
		Name:   "artist_name",
		Images: images,
	}

	images = []*model.Image{
		{
			Height: &imageSizeOne,
			Width:  &imageSizeOne,
			Url:    "url",
		},
		{
			Height: &imageSizeTwo,
			Width:  &imageSizeTwo,
			Url:    "url",
		},
	}
	imagesBson = bson.A{
		bson.D{
			{Key: "height", Value: 2},
			{Key: "width", Value: 2},
			{Key: "url", Value: "url"},
		},
		bson.D{
			{Key: "height", Value: 1},
			{Key: "width", Value: 1},
			{Key: "url", Value: "url"},
		},
	}
)

func TestFind(t *testing.T) {
	testCases := []struct {
		name            string
		artworkID       string
		dbResponse      []bson.D
		expectedArtwork *model.Artwork
		expectedError   error
	}{
		{
			name:            "Invalid artworkID",
			artworkID:       "invalid_ID",
			dbResponse:      []bson.D{},
			expectedArtwork: nil,
			expectedError:   primitive.ErrInvalidHex,
		},
		{
			name:      "no artwork found",
			artworkID: artworkID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artwork", mtest.FirstBatch),
			},
			expectedArtwork: nil,
			expectedError:   mongo.ErrNoDocuments,
		},
		{
			name:      "artwork found",
			artworkID: artworkID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artwork", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: artworkID},
					{Key: "title", Value: "artwork_title"},
					{Key: "images", Value: imagesBson},
					{Key: "artist", Value: bson.D{
						{Key: "_id", Value: artistID},
						{Key: "name", Value: "artist_name"},
						{Key: "images", Value: imagesBson},
					}},
				}),
			},
			expectedArtwork: &model.Artwork{
				ID:     artworkObjectID,
				Title:  "artwork_title",
				Images: images,
				Artist: artist,
			},
			expectedError: nil,
		},
		{
			name:      "find returns error",
			artworkID: artworkID,
			dbResponse: []bson.D{
				MongoFailResponse,
			},
			expectedArtwork: nil,
			expectedError:   ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			artwork, err := store.Find(context.Background(), tc.artworkID)
			require.Equal(mt, tc.expectedArtwork, artwork)
			require.Equal(mt, tc.expectedError, err)
		})
	}

}

func TestFindMany(t *testing.T) {
	testCases := []struct {
		name            string
		dbResponse      []bson.D
		expectedArtwork []*model.Artwork
		expectedError   error
	}{
		{
			name: "no artwork found",
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artworks", mtest.FirstBatch),
			},
			expectedArtwork: []*model.Artwork{},
			expectedError:   nil,
		},
		{
			name: "artworks found",
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artworks", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: artworkID},
						{Key: "title", Value: "title_one"},
						{Key: "images", Value: imagesBson},
						{Key: "artist", Value: artist},
					},
					bson.D{
						{Key: "_id", Value: artworkID},
						{Key: "title", Value: "title_two"},
						{Key: "images", Value: imagesBson},
						{Key: "artist", Value: artist},
					},
				),
			},
			expectedArtwork: []*model.Artwork{
				{
					ID:     artworkObjectID,
					Title:  "title_one",
					Images: images,
					Artist: artist,
				},
				{
					ID:     artworkObjectID,
					Title:  "title_two",
					Images: images,
					Artist: artist,
				},
			},
			expectedError: nil,
		},
		{
			name: "find returns error",
			dbResponse: []bson.D{
				MongoFailResponse,
			},
			expectedArtwork: nil,
			expectedError:   ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			artwork, err := store.FindMany(context.Background())
			require.Equal(mt, tc.expectedArtwork, artwork)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestInsertMany(t *testing.T) {
	artworkObjectIDTwo, _ := primitive.ObjectIDFromHex("60e0850266d6c13d7b599b6b")
	artworks := []*model.Artwork{
		{
			ID:     artworkObjectID,
			Title:  "title_one",
			Images: images,
			Artist: artist,
		},
		{
			ID:     artworkObjectIDTwo,
			Title:  "title_two",
			Images: images,
			Artist: artist,
		},
	}

	testCases := []struct {
		name          string
		artworks      []*model.Artwork
		dbResponse    []bson.D
		expectedError error
	}{
		{
			name:     "insert artists",
			artworks: artworks,
			dbResponse: []bson.D{
				mtest.CreateSuccessResponse(),
			},
			expectedError: nil,
		},
		{
			name:     "insert many fails with an error",
			artworks: artworks,
			dbResponse: []bson.D{
				MongoFailResponse,
			},
			expectedError: ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			err := store.InsertMany(context.Background(), tc.artworks)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}
