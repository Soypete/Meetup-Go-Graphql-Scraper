package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kollalabs/sdk-go/kc"
)

type payloadql struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

type ProNetworkByUrlname struct {
	Data struct {
		ProNetwork struct {
			GroupsSearch GroupsSearch `json:"groupsSearch,omitempty"`
			EventsSearch EventsSearch `json:"eventsSearch,omitempty"`
		} `json:"proNetworkByUrlname"`
	} `json:"data"`
}
type GroupsSearch struct {
	Count    int `json:"count"`
	PageInfo struct {
		HasNextPage bool   `json:"hasNextPage"`
		StartCursor string `json:"startCursor"`
		EndCursor   string `json:"endCursor"`
	} `json:"pageInfo"`
	Edges []struct {
		Node struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"node"`
	} `json:"edges"`
}

type EventsSearch struct {
	Count    int `json:"count"`
	PageInfo struct {
		HasNextPage bool   `json:"hasNextPage"`
		StartCursor string `json:"startCursor"`
		EndCursor   string `json:"endCursor"`
	} `json:"pageInfo"`
	Edges []struct {
		Node struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Group struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"group"`
			DateTime string `json:"dateTime"`
		} `json:"node"`
	} `json:"edges"`
}

func getBearerToken() (string, error) {
	// Get api key from environment variable

	// TODO(soypete 03-27-2023): rename to kolla key
	apiKey := os.Getenv("KOLLA_KEY")
	ctx := context.Background()

	if apiKey == "" {
		return "", fmt.Errorf("no kolla key provided.")
	}
	// Create a new client
	kolla, err := kc.New(apiKey)
	if err != nil {
		return "", fmt.Errorf("unable to load kolla connect client: %s", err)
	}
	// Get consumer token
	// TODO(soypete 03-27-2023): update names to correspond wiht kolla api
	if err != nil {
		return "", fmt.Errorf("unable to load consumer token: %s", err)
	}
	creds, err := kolla.Credentials(ctx, os.Getenv("CONNECTOR_ID"), os.Getenv("CONSUMER_ID"))
	if err != nil {
		return "", fmt.Errorf("unable to load consumer token: %s", err)
	}

	return fmt.Sprintf("Bearer %s", creds.Token), nil
}

func getListOfGroups(bearerToken string) {
	query := `query ($urlname: String!) { 
		proNetworkByUrlname(urlname: $urlname) { 
			groupsSearch(input: {first: 3}) {
      count
      pageInfo {
				hasNextPage
				startCursor
        endCursor
      }
      edges {
        node {
          id
          name
				} 
			} 
		} 
	} 
}
  `
	variables := `{"urlname":"forge-utah"}`
	p := payloadql{
		Query:     query,
		Variables: variables,
	}
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://api.meetup.com/gql", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	// set header fields
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// run it and capture the response
	var respData ProNetworkByUrlname

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatal(err)
	}
	meetupGroupIds := respData.Data.ProNetwork.GroupsSearch.Edges
	fmt.Print(meetupGroupIds)
}

func getListOfEvents(bearerToken string) {
	query := `query ($urlname: String!) { 
		proNetworkByUrlname(urlname: $urlname) { 
			eventsSearch(input: {first: 3}) {
      count
      pageInfo {
				hasNextPage
				startCursor
        endCursor
      }
      edges {
        node {
          id
         	title
					group {
						id
						name
					}
					dateTime
				} 
			} 
		} 
	} 
}
  `
	variables := `{"urlname":"forge-utah"}`
	p := payloadql{
		Query:     query,
		Variables: variables,
	}
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://api.meetup.com/gql", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	// set header fields
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// run it and capture the response
	var respData ProNetworkByUrlname

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("number of events %d\n", len(respData.Data.ProNetwork.EventsSearch.Edges))
}
func main() {
	bearerToken, err := getBearerToken()
	if err != nil {
		log.Fatal(err)
	}
	getListOfGroups(bearerToken)
	getListOfEvents(bearerToken)
}
