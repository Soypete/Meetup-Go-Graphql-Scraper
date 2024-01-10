package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/soypete/Metup-Go-Graphql-Scraper/auth"
	"github.com/soypete/Metup-Go-Graphql-Scraper/meetup"
)

type Config struct {
	Key         string `json:"kolla_key"`
	ConsumerID  string `json:"consumer_id"`
	ConnectorID string `json:"connector_id"`
	ProAccount  string `json:"pro_account"`
}

func parseConfig(file string) (Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open config file: %s", err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	config := Config{}
	if err = jsonParser.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("unable to parse config file: %s", err)
	}
	return config, nil
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "config.json", "config file path")
	flag.Parse()

	parsedConfig, err := parseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	kc := auth.Setup(parsedConfig.Key, parsedConfig.ConsumerID, parsedConfig.ConnectorID)
	bearerToken, err := kc.GetBearerToken()
	if err != nil {
		log.Fatal(err)
	}
	meetupClient := meetup.Setup(bearerToken, parsedConfig.ProAccount)

	// TODO(soypete): create csv with relevant data.
	// for list of data in README.md
	err = meetupClient.BuildDataSet()
	if err != nil {
		log.Fatal(err)
	}
}
