package config

import (
	"fmt"
	"io/ioutil"
)

//easyjson:json
type Database struct {
	SchemaFile      string
	TimestampFormat string
}

//easyjson:json
type API struct {
	TimestampFormat string
}

//easyjson:json
type Config struct {
	Database Database
	API      API
}

// Instance is a singleton of configuration
var Instance = Config{}

func init() {
	configBytes, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := Instance.UnmarshalJSON(configBytes); err != nil {
		fmt.Println(err)
		return
	}
}
