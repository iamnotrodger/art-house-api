package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseArtists(t *testing.T) {
	artists, err := parseArtists("./data/artists.json")
	require.Equal(t, nil, err)
	t.Log(artists)
}

func TestParseArtworks(t *testing.T) {
	artworks, err := parseArtworks("./data/artworks.json")
	require.Equal(t, nil, err)
	t.Log(artworks)
}

func TestParseExhibitions(t *testing.T) {
	exhibitions, err := parseExhibitions("./data/exhibitions.json")
	require.Equal(t, nil, err)
	t.Log(exhibitions)
}
