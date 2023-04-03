// meetup module contains all the logic for making our
// data requests to Meetup.com's grapql-api.
// The api docs can be found here: https://www.meetup.com/api/guide/#graphQl-guide.
package meetup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client connects all information for connecting to the Meetup API.
type Client struct {
	ql          *http.Client
	bearerToken string //TODO(soypete): refresh
	url         string
}

func Setup(bearerToken string) Client {
	return Client{
		ql: &http.Client{
			Timeout: 10 * time.Second,
		},
		bearerToken: bearerToken,
		url:         "https://api.meetup.com/gql",
	}
}

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
	Edges []Edge `json:"edges"`
}

type Edge struct {
	Node Node `json:"node"`
}
type Node struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	Name  string `json:"name,omitempty"`
	Group struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"group,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
}

type EventsSearch struct {
	Count    int `json:"count"`
	PageInfo struct {
		HasNextPage bool   `json:"hasNextPage"`
		StartCursor string `json:"startCursor"`
		EndCursor   string `json:"endCursor"`
	} `json:"pageInfo"`
	Edges []Edge `json:"edges"`
}

func (c Client) GetListOfGroups() ProNetworkByUrlname {
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

	body, err := c.sendRequest(p)
	if err != nil {
		log.Fatal(err)
	}
	var respData ProNetworkByUrlname
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatal(err)
	}
	return respData
}

func (c Client) GetListOfEvents() ProNetworkByUrlname {
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
	body, err := c.sendRequest(p)
	if err != nil {
		log.Fatal(err)
	}
	// run it and capture the response
	var respData ProNetworkByUrlname
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatal(err)
	}
	return respData
}

func (c Client) sendRequest(ql payloadql) (resp []byte, err error) {
	b, err := json.Marshal(ql)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal graphql payload, %w", err)
	}
	req, err := http.NewRequest("POST", c.url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create new http.Request, %w", err)
	}
	// set header fields
	req.Header.Add("Authorization", c.bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.ql.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed meetup.com api request, %w", err)
	}
	// TODO(soypete): check payload status
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload body, %w", err)
	}
	return body, nil
}
