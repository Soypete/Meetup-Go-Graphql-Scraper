package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

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
	var analyticsFunc string
	flag.StringVar(&configPath, "config-path", "config.json", "config file path")
	flag.StringVar(&analyticsFunc, "analytics", "[groups, eventRSVP]", "function to run for analytics")
	flag.Parse()

	// parse config file
	parsedConfig, err := parseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	// get meetup bearer token
	kc := auth.Setup(parsedConfig.Key, parsedConfig.ConsumerID, parsedConfig.ConnectorID)
	bearerToken, err := kc.GetBearerToken()
	if err != nil {
		log.Fatal(err)
	}
	meetupClient := meetup.Setup(bearerToken, parsedConfig.ProAccount)
	listFuncs, err := meetupClient.GetAnalyticsFunc(analyticsFunc)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	// run analytics function
	wg.Add(len(listFuncs))
	for _, f := range listFuncs {
		go func(f func()) {
			f()
			defer wg.Done()
		}(f)
	}
	wg.Wait()
}
