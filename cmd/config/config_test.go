package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	assert.Equal(t, Global.Port, defaultPort)
	assert.Equal(t, Global.MongoDBName, defaultMongoDBName)
	assert.Equal(t, Global.MongoURI, defaultMongoURI)

	assert.Equal(t, Global.RedisAddr, defaultRedisAddr)
	assert.Equal(t, Global.RedisPassword, defaultRedisPassword)
	assert.Equal(t, Global.RedisDb, defaultRedisDb)

	assert.Equal(t, Global.ArtworkLimit, defaultArtworkLimit)
	assert.Equal(t, Global.ArtworkLimitMin, defaultArtworkLimitMin)
	assert.Equal(t, Global.ArtworkLimitMax, defaultArtworkLimitMax)

	assert.Equal(t, Global.ArtistLimit, defaultArtistLimit)
	assert.Equal(t, Global.ArtistLimitMin, defaultArtistLimitMin)
	assert.Equal(t, Global.ArtistLimitMax, defaultArtistLimitMax)

	assert.Equal(t, Global.ExhibitionLimit, defaultExhibitionLimit)
	assert.Equal(t, Global.ExhibitionLimitMin, defaultExhibitionLimitMin)
	assert.Equal(t, Global.ExhibitionLimitMax, defaultExhibitionLimitMax)
}
