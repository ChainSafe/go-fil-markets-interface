package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	NodeIP          string
	NodeAuthToken   string
	MarketIP        string
	MarketAuthToken string
}

var Market Config

func Load(configFile string) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Unable to load config from file : %v", err))
	}

	err = json.Unmarshal(file, &Market)
	if err != nil {
		panic(fmt.Sprintf("Unable to unmarshall config from file : %v", err))
	}
}
