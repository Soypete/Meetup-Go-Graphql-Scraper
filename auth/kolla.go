package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/kollalabs/sdk-go/kc"
)

func GetBearerToken() (string, error) {
	// Get api key from environment variable

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
	if err != nil {
		return "", fmt.Errorf("unable to load consumer token: %s", err)
	}
	creds, err := kolla.Credentials(ctx, os.Getenv("CONNECTOR_ID"), os.Getenv("CONSUMER_ID"))
	if err != nil {
		return "", fmt.Errorf("unable to load consumer token: %s", err)
	}

	return fmt.Sprintf("Bearer %s", creds.Token), nil
}
