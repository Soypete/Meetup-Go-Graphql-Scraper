package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kollalabs/sdk-go/kc"
	"github.com/machinebox/graphql"
)

func getBearerToken() (string, error) {
	// Get api key from environment variable
	apiKey := os.Getenv("MEETUP_API_KEY")
	ctx := context.Background()

	if apiKey == "" {
		return "", fmt.Errorf("no meetup api key provided")
	}
	// Create a new client
	kolla, err := kc.New(apiKey)
	if err != nil {
		return "", fmt.Errorf("unable to load kolla connect client: %s", err)
	}
	// Get consumer token
	consumerToken, err := kolla.ConsumerToken(ctx, os.Getenv("MEETUP_TOKEN"), os.Getenv("MEETUP_USERNAME"))
	if err != nil {
		return "", fmt.Errorf("unable to load consumer token: %s", err)
	}
	return fmt.Sprintf("Bearer %s", consumerToken), nil
}

func main() {
	bearerToken, err := getBearerToken()
	if err != nil {
		log.Fatal(err)
	}
	client := graphql.NewClient("https://api.meetup.com/gql")

	// make a request
	req := graphql.NewRequest(`
    query {
        self {
           id 
           name 
        }
    }
`)

	// set header fields
	req.Header.Add("Authorization", bearerToken)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	type Data struct {
		Self struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"self"`
	}
	var respData Data
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	fmt.Println(respData)
}
