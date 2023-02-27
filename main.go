package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/kollalabs/sdk-go/kc"
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
	ctx := context.WithValue(context.Background(), "Authorization", bearerToken)

	// make graphql query
	// query := `
	// {
	// 	checkIfGroupUrlnameValid(urlname: "utahgophers") {
	// 		isValid
	// 		urlname
	// 		error
	// 	}
	// }
	// `
	query := `
{
	self {
		memberships
		tickets
		id
		name
	}
}
`

	// execute graphql request
	params := graphql.Params{RequestString: query, Context: ctx}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("failed to unmarshal graphql respons : %+v", err)
	}

	fmt.Printf("%s \n", rJSON)
}
