package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	assert.Equal(t, Global.Port, defaultPort)
	assert.Equal(t, Global.MongoDBName, defaultMongoDBName)
	assert.Equal(t, Global.MongoURI, defaultMongoURI)
}
