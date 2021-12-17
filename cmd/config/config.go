package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	defaultPort        = ":8080"
	defaultMongoURI    = "mongodb://localhost:27017"
	defaultMongoDBName = "art-house"

	defaultArtworkLimit    = 15
	defaultArtworkLimitMin = 1
	defaultArtworkLimitMax = 100

	defaultArtistLimit    = 15
	defaultArtistLimitMin = 1
	defaultArtistLimitMax = 100

	defaultExhibitionLimit    = 15
	defaultExhibitionLimitMin = 1
	defaultExhibitionLimitMax = 100
)

type Spec struct {
	Port               string `mapstructure:"port"`
	MongoURI           string `mapstructure:"mongo_uri"`
	MongoDBName        string `mapstructure:"mongo_db_name"`
	ArtworkLimit       int64  `mapstructure:"artwork_limit"`
	ArtworkLimitMin    int64  `mapstructure:"artwork_limit_min"`
	ArtworkLimitMax    int64  `mapstructure:"artwork_limit_max"`
	ArtistLimit        int64  `mapstructure:"artist_limit"`
	ArtistLimitMin     int64  `mapstructure:"artist_limit_min"`
	ArtistLimitMax     int64  `mapstructure:"artist_limit_max"`
	ExhibitionLimit    int64  `mapstructure:"exhibition_limit"`
	ExhibitionLimitMin int64  `mapstructure:"exhibition_limit_min"`
	ExhibitionLimitMax int64  `mapstructure:"exhibition_limit_max"`
}

var Global = Spec{
	Port:               defaultPort,
	MongoURI:           defaultMongoURI,
	MongoDBName:        defaultMongoDBName,
	ArtworkLimit:       defaultArtworkLimit,
	ArtworkLimitMin:    defaultArtworkLimitMin,
	ArtworkLimitMax:    defaultArtworkLimitMax,
	ArtistLimit:        defaultArtistLimit,
	ArtistLimitMin:     defaultArtistLimitMin,
	ArtistLimitMax:     defaultArtistLimitMax,
	ExhibitionLimit:    defaultExhibitionLimit,
	ExhibitionLimitMin: defaultExhibitionLimitMin,
	ExhibitionLimitMax: defaultExhibitionLimitMax,
}

func init() {
	v := viper.New()
	v.SetConfigFile(".env")
	v.ReadInConfig()
	v.AutomaticEnv()

	if err := v.Unmarshal(&Global); err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config %s", err))
	}
}
