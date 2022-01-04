package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortImages(t *testing.T) {
	var (
		size500  float64 = 500
		size1080 float64 = 1080
		size1920 float64 = 1920
	)

	tests := []struct {
		name           string
		images         []*Image
		expectedImages []*Image
	}{
		{
			name:           "sorting nil images",
			images:         nil,
			expectedImages: nil,
		},
		{
			name: "sorting images",
			images: []*Image{
				{
					Height: &size1920,
					Width:  &size1920,
				},
				{
					Height: &size1080,
					Width:  &size1080,
				},
				{
					Height: &size500,
					Width:  &size500,
				},
				{
					Height: &size500,
					Width:  &size500,
				},
				{
					Height: nil,
					Width:  nil,
				},
			},
			expectedImages: []*Image{
				{
					Height: &size500,
					Width:  &size500,
				},
				{
					Height: &size500,
					Width:  &size500,
				},
				{
					Height: &size1080,
					Width:  &size1080,
				},
				{
					Height: &size1920,
					Width:  &size1920,
				},
				{
					Height: nil,
					Width:  nil,
				},
			},
		},
		{
			name: "already sorted images",
			images: []*Image{
				{
					Height: &size500,
					Width:  &size500,
				},
				{
					Height: &size1080,
					Width:  &size1080,
				},
				{
					Height: &size1920,
					Width:  &size1920,
				},
				{
					Height: nil,
					Width:  nil,
				},
			},
			expectedImages: []*Image{
				{
					Height: &size500,
					Width:  &size500,
				},
				{
					Height: &size1080,
					Width:  &size1080,
				},
				{
					Height: &size1920,
					Width:  &size1920,
				},
				{
					Height: nil,
					Width:  nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SortImages(test.images)
			require.Equal(t, test.images, test.expectedImages)
		})
	}
}
