package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	defaultPort               = 8080
	defaultMongoURI           = "mongodb://localhost:27017"
	defaultMongoDBName        = "art-house"
	defaultRedisAddr          = "localhost:6379"
	defaultRedisPassword      = ""
	defaultRedisDb            = 0
	defaultArtworkLimit       = int64(15)
	defaultArtworkLimitMin    = int64(1)
	defaultArtworkLimitMax    = int64(100)
	defaultArtistLimit        = int64(15)
	defaultArtistLimitMin     = int64(1)
	defaultArtistLimitMax     = int64(100)
	defaultExhibitionLimit    = int64(15)
	defaultExhibitionLimitMin = int64(1)
	defaultExhibitionLimitMax = int64(100)
)

type Spec struct {
	Port               int    `mapstructure:"port"`
	MongoURI           string `mapstructure:"mongo_uri"`
	MongoDBName        string `mapstructure:"mongo_db_name"`
	RedisAddr          string `mapstructure:"redis_addr"`
	RedisPassword      string `mapstructure:"redis_password"`
	RedisDb            int    `mapstructure:"redis_db"`
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
	RedisAddr:          defaultRedisAddr,
	RedisPassword:      defaultRedisPassword,
	RedisDb:            defaultRedisDb,
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

func LoadConfig() {
	v := viper.New()
	v.SetConfigFile(".env")
	v.ReadInConfig()
	v.AutomaticEnv()

	setDefaults(v, Global)

	if err := v.Unmarshal(&Global); err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config %s", err))
	}
}

func setDefaults(v *viper.Viper, i interface{}) {
	values := map[string]interface{}{}
	if err := mapstructure.Decode(i, &values); err != nil {
		panic(err)
	}
	for key, defaultValue := range values {
		v.SetDefault(key, defaultValue)
	}
}
