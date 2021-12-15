package exhibition

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
	exhibitID          = "60e0850266d6c13d7b599b6d"
	artworkObjectID, _ = primitive.ObjectIDFromHex(artworkID)
	artistObjectID, _  = primitive.ObjectIDFromHex(artistID)
	exhibitObjectID, _ = primitive.ObjectIDFromHex(exhibitID)
	imageSizeOne       = 1.0
	imageSizeTwo       = 2.0

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
		exhibitID       string
		dbResponse      []bson.D
		expectedExhibit *model.Exhibition
		expectedError   error
	}{
		{
			name:            "Invalid artworkID",
			exhibitID:       "invalid_ID",
			dbResponse:      []bson.D{},
			expectedExhibit: nil,
			expectedError:   primitive.ErrInvalidHex,
		},
		{
			name:      "no exhibition found",
			exhibitID: exhibitID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.exhibitions", mtest.FirstBatch),
			},
			expectedExhibit: nil,
			expectedError:   mongo.ErrNoDocuments,
		},
		{
			name:      "exhibition found",
			exhibitID: exhibitID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.exhibitions", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: exhibitID},
					{Key: "name", Value: "exhibition_name"},
					{Key: "images", Value: imagesBson},
				}),
			},
			expectedExhibit: &model.Exhibition{
				ID:     exhibitObjectID,
				Name:   "exhibition_name",
				Images: images,
			},
			expectedError: nil,
		},
		{
			name:      "find returns error",
			exhibitID: exhibitID,
			dbResponse: []bson.D{
				MongoFailResponse,
			},
			expectedExhibit: nil,
			expectedError:   ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			exhibit, err := store.Find(context.Background(), tc.exhibitID)
			require.Equal(mt, tc.expectedExhibit, exhibit)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestFindMany(t *testing.T) {
	testCases := []struct {
		name                string
		dbResponse          []bson.D
		expectedExhibitions []*model.Exhibition
		expectedError       error
	}{
		{
			name: "no exhibitions found",
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.exhibitions", mtest.FirstBatch),
			},
			expectedExhibitions: []*model.Exhibition{},
			expectedError:       nil,
		},
		{
			name: "exhibitions found",
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.exhibitions", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: exhibitID},
						{Key: "name", Value: "exhibit_name_1"},
						{Key: "images", Value: imagesBson},
					},
					bson.D{
						{Key: "_id", Value: exhibitID},
						{Key: "name", Value: "exhibit_name_2"},
						{Key: "images", Value: imagesBson},
					},
				),
			},
			expectedExhibitions: []*model.Exhibition{
				{
					ID:     exhibitObjectID,
					Name:   "exhibit_name_1",
					Images: images,
				},
				{
					ID:     exhibitObjectID,
					Name:   "exhibit_name_2",
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
			expectedExhibitions: nil,
			expectedError:       ErrMongoCommandError,
		},
	}

	m := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer m.Close()
	for _, tc := range testCases {
		m.Run(tc.name, func(mt *mtest.T) {
			mt.AddMockResponses(tc.dbResponse...)

			store := NewStore(mt.DB)
			exhibitions, err := store.FindMany(context.Background(), bson.D{})
			require.Equal(mt, tc.expectedExhibitions, exhibitions)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestFindArtworks(t *testing.T) {
	testCases := []struct {
		name             string
		exhibitionID     string
		dbResponse       []bson.D
		expectedArtworks []*model.Artwork
		expectedError    error
	}{
		{
			name:             "invalid exhibitID",
			exhibitionID:     "invalid_ID",
			dbResponse:       []bson.D{},
			expectedArtworks: nil,
			expectedError:    primitive.ErrInvalidHex,
		},
		{
			name:         "no exhibit's artwork found",
			exhibitionID: exhibitID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(1, "art-house.artworks", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: exhibitID},
					},
				),
			},
			expectedArtworks: nil,
			expectedError:    nil,
		},
		{
			name:         "exhibit's artwork found",
			exhibitionID: exhibitID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artworks", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: exhibitID},
						{Key: "artworks", Value: bson.A{
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
						},
						},
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
			name:         "find returns error",
			exhibitionID: exhibitID,
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
			artworks, err := store.FindArtworks(context.Background(), tc.exhibitionID)
			require.Equal(mt, tc.expectedArtworks, artworks)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestFindArtists(t *testing.T) {
	testCases := []struct {
		name            string
		exhibitionID    string
		dbResponse      []bson.D
		expectedArtists []*model.Artist
		expectedError   error
	}{
		{
			name:            "invalid exhibitID",
			exhibitionID:    "invalid_ID",
			dbResponse:      []bson.D{},
			expectedArtists: nil,
			expectedError:   primitive.ErrInvalidHex,
		},
		{
			name:         "no exhibit's artwork found",
			exhibitionID: exhibitID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(1, "art-house.artworks", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: exhibitID},
					},
				),
			},
			expectedArtists: nil,
			expectedError:   nil,
		},
		{
			name:         "exhibit's artwork found",
			exhibitionID: exhibitID,
			dbResponse: []bson.D{
				mtest.CreateCursorResponse(0, "art-house.artworks", mtest.FirstBatch,
					bson.D{
						{Key: "_id", Value: exhibitID},
						{Key: "artists", Value: bson.A{
							bson.D{
								{Key: "_id", Value: artistID},
								{Key: "name", Value: "name_1"},
								{Key: "images", Value: imagesBson},
							},
							bson.D{
								{Key: "_id", Value: artistID},
								{Key: "name", Value: "name_2"},
								{Key: "images", Value: imagesBson},
							},
						},
						},
					},
				),
			},
			expectedArtists: []*model.Artist{
				{
					ID:     artistObjectID,
					Name:   "name_1",
					Images: images,
				},
				{
					ID:     artistObjectID,
					Name:   "name_2",
					Images: images,
				},
			},
			expectedError: nil,
		},
		{
			name:         "find returns error",
			exhibitionID: exhibitID,
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
			artists, err := store.FindArtists(context.Background(), tc.exhibitionID)
			require.Equal(mt, tc.expectedArtists, artists)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}

func TestInsertMany(t *testing.T) {
	exhibitions := []*model.Exhibition{
		{
			ID:     exhibitObjectID,
			Name:   "name_one",
			Images: images,
			Artists: []*model.Artist{
				{
					ID: artistObjectID,
				},
			},
			Artworks: []*model.Artwork{
				{
					ID: artworkObjectID,
				},
			},
		},
		{
			ID:     exhibitObjectID,
			Name:   "name_two",
			Images: images,
			Artists: []*model.Artist{
				{
					ID: artistObjectID,
				},
			},
			Artworks: []*model.Artwork{
				{
					ID: artworkObjectID,
				},
			},
		},
	}

	testCases := []struct {
		name          string
		exhibitions   []*model.Exhibition
		dbResponse    []bson.D
		expectedError error
	}{
		{
			name:        "insert artists",
			exhibitions: exhibitions,
			dbResponse: []bson.D{
				mtest.CreateSuccessResponse(),
			},
			expectedError: nil,
		},
		{
			name:        "insert many fails with an error",
			exhibitions: exhibitions,
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
			err := store.InsertMany(context.Background(), tc.exhibitions)
			require.Equal(mt, tc.expectedError, err)
		})
	}
}
