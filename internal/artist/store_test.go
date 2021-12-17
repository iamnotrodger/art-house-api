package artist

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

	artistID          = "60e0850266d6c13d7b599b69"
	artistObjectID, _ = primitive.ObjectIDFromHex(artistID)
	imageSizeOne      = 1.0
	imageSizeTwo      = 2.0

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
		name           string
		artistID       string
		dbResponse     []bson.D
		expectedArtist *model.Artist
		expectedError  error
	}{
		{
			name:           "Invalid artworkID",
			artistID:       "invalid_ID",
			dbResponse:     []bson.D{},
			expectedArtist: nil,
			expectedError:  primitive.ErrInvalidHex,
		},
		{
			name:     "no artist found",
			artistID: artistID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artists", mtest.FirstBatch),
			},
			expectedArtist: nil,
			expectedError:  mongo.ErrNoDocuments,
		},
		{
			name:     "artist found",
			artistID: artistID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artists", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: artistID},
					{Key: "name", Value: "artist_name"},
					{Key: "images", Value: imagesBson},
				}),
			},
			expectedArtist: &model.Artist{
				ID:     artistObjectID,
				Name:   "artist_name",
				Images: images,
			},
			expectedError: nil,
		},
		{
			name:     "find returns error",
			artistID: artistID,
			dbResponse: []bson.D{
				MongoFailResponse,
			},
			expectedArtist: nil,
			expectedError:  ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			artist, err := store.Find(context.Background(), tc.artistID)
			require.Equal(mt, tc.expectedArtist, artist)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestFindMany(t *testing.T) {
	testCases := []struct {
		name            string
		dbResponse      []bson.D
		expectedArtists []*model.Artist
		expectedError   error
	}{
		{
			name: "no artists found",
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artists", mtest.FirstBatch),
			},
			expectedArtists: []*model.Artist{},
			expectedError:   nil,
		},
		{
			name: "artists found",
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artists", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: artistID},
						{Key: "name", Value: "artist_name_1"},
						{Key: "images", Value: imagesBson},
					},
					bson.D{
						{Key: "_id", Value: artistID},
						{Key: "name", Value: "artist_name_2"},
						{Key: "images", Value: imagesBson},
					},
				),
			},
			expectedArtists: []*model.Artist{
				{
					ID:     artistObjectID,
					Name:   "artist_name_1",
					Images: images,
				},
				{
					ID:     artistObjectID,
					Name:   "artist_name_2",
					Images: images,
				},
			},
			expectedError: nil,
		},
		{
			name: "find returns error",
			dbResponse: []bson.D{
				MongoFailResponse,
			},
			expectedArtists: nil,
			expectedError:   ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			artist, err := store.FindMany(context.Background())
			require.Equal(mt, tc.expectedArtists, artist)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestFindArtworks(t *testing.T) {
	artworkID := "60e0850266d6c13d7b599b69"
	artworkObjectID, _ := primitive.ObjectIDFromHex(artworkID)

	testCases := []struct {
		name             string
		artistID         string
		dbResponse       []bson.D
		expectedArtworks []*model.Artwork
		expectedError    error
	}{
		{
			name:             "invalid artworkID",
			artistID:         "invalid_ID",
			dbResponse:       []bson.D{},
			expectedArtworks: nil,
			expectedError:    primitive.ErrInvalidHex,
		},
		{
			name:     "no artist's artwork found",
			artistID: artistID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artworks", mtest.FirstBatch),
			},
			expectedArtworks: []*model.Artwork{},
			expectedError:    nil,
		},
		{
			name:     "artist's artwork found",
			artistID: artistID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artworks", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: artworkID},
						{Key: "title", Value: "title_1"},
						{Key: "images", Value: imagesBson},
						{Key: "year", Value: 1},
						{Key: "description", Value: "description"},
					},
					bson.D{
						{Key: "_id", Value: artworkID},
						{Key: "title", Value: "title_2"},
						{Key: "images", Value: imagesBson},
						{Key: "year", Value: 1},
						{Key: "description", Value: "description"},
					},
				),
			},
			expectedArtworks: []*model.Artwork{
				{
					ID:          artworkObjectID,
					Title:       "title_1",
					Images:      images,
					Year:        1,
					Description: "description",
				},
				{
					ID:          artworkObjectID,
					Title:       "title_2",
					Images:      images,
					Year:        1,
					Description: "description",
				},
			},
			expectedError: nil,
		},
		{
			name:     "find returns error",
			artistID: artworkID,
			dbResponse: []bson.D{
				MongoFailResponse,
			},
			expectedArtworks: nil,
			expectedError:    ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			artworks, err := store.FindArtworks(context.Background(), tc.artistID)
			require.Equal(mt, tc.expectedArtworks, artworks)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestInsertMany(t *testing.T) {
	artistObjectIDTwo, _ := primitive.ObjectIDFromHex("60e0850266d6c13d7b599b6a")
	artists := []*model.Artist{
		{
			ID:     artistObjectID,
			Name:   "name_one",
			Images: images,
		},
		{
			ID:     artistObjectIDTwo,
			Name:   "name_two",
			Images: images,
		},
	}

	testCases := []struct {
		name          string
		artists       []*model.Artist
		dbResponse    []bson.D
		expectedError error
	}{
		{
			name:    "insert artists",
			artists: artists,
			dbResponse: []bson.D{
				mtest.CreateSuccessResponse(),
			},
			expectedError: nil,
		},
		{
			name:    "insert many fails with an error",
			artists: artists,
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
			err := store.InsertMany(context.Background(), tc.artists)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}
