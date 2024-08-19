package auth

import (
	"context"

	"github.com/coreos/go-oidc"
)

var (
	clientID     = "your-client-id"
	clientSecret = "your-client-secret"
	providerURL  = "https://accounts.google.com"
)

func InitOIDC() (*oidc.Provider, *oidc.IDTokenVerifier, error) {
	provider, err := oidc.NewProvider(context.Background(), providerURL)
	if err != nil {
		return nil, nil, err
	}

	config := &oidc.Config{
		ClientID: clientID,
	}

	verifier := provider.Verifier(config)
	return provider, verifier, nil
}
