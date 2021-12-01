package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	defaultPort        = ":8080"
	defaultMongoURI    = "mongodb://localhost:27017"
	defaultMongoDBName = "art-house"
)

type Spec struct {
	Port        string `mapstructure:"port"`
	MongoURI    string `mapstructure:"mongo_uri"`
	MongoDBName string `mapstructure:"mongo_db_name"`
}

var Global = Spec{
	Port:        defaultPort,
	MongoURI:    defaultMongoURI,
	MongoDBName: defaultMongoDBName,
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
