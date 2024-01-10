package auth

import (
	"context"
	"fmt"

	kolla "github.com/kollalabs/sdk-go/kc"
)

type KollaConnect struct {
	key         string
	consumerID  string
	connectorID string
}

func Setup(key, consumerID, connectorID string) *KollaConnect {
	return &KollaConnect{
		key:         key,
		consumerID:  consumerID,
		connectorID: connectorID,
	}
}

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
