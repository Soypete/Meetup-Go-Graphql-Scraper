package auth

import (
	"context"
	"fmt"

	kolla "github.com/kollalabs/sdk-go/kc"
)

// KollaConnect is a struct that holds the key, consumerID and connectorID.
type KollaConnect struct {
	key         string
	consumerID  string
	connectorID string
}

// Setup is a function that returns a new KollaConnect struct.
func Setup(key, consumerID, connectorID string) *KollaConnect {
	return &KollaConnect{
		key:         key,
		consumerID:  consumerID,
		connectorID: connectorID,
	}
}

// GetBearerToken is a function that uses the kolla connect client to get the meetup api
// oauth credentials.
func (kc *KollaConnect) GetBearerToken() (string, error) {
	// Get api key from environment variable

	ctx := context.Background()

	// Create a new client
	kolla, err := kolla.New(kc.key)
	if err != nil {
		return "", fmt.Errorf("unable to load kolla connect client: %s", err)
	}
	creds, err := kolla.Credentials(ctx, kc.connectorID, kc.consumerID)
	if err != nil {
		return "", fmt.Errorf("unable to load consumer token: %s", err)
	}

	return fmt.Sprintf("Bearer %s", creds.Token), nil
}
