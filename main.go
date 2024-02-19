/*
Meetup-Go-Graphql-Scraper is a tool to scrape data from Meetup.com using their GraphQL API.
It allows users to specify which analytics function to run and will output the results to a file.
The current analytics functions are:
- groups: list all the groups associate to a pro account
- eventRSVP: list all the RSVPs for a events in a pro account
- groupAnalytics: export of group analytics data as provided by Meetup.com. So far this has been empty for all groups.

The Kolla integration tool us used to authenticate with Meetup.com and get a bearer token for meetup.com's OAuth2 API.
To use this tool, you will need to have a Kolla account and a Meetup.com pro account, as well as a config file with the following fields:

```
{
  "pro_account": "go",
  "kolla_key": {kolla.secret},
  "connector_id": {kolla.account},
  "consumer_id": {kolla.key}
}
```

usage:
```
go run meetup-go-grapghl-scraper -config-path=config.json -analytics=[groups, eventRSVP]
```

the flags are:
- config-path: the path to the config file
- analytics: the function to run for analytics

*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/soypete/Meetup-Go-Graphql-Scraper/auth"
	"github.com/soypete/Meetup-Go-Graphql-Scraper/meetup"
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
