package auth

import (
	"log"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

// ConfigureOIDCProvider sets up the OIDC provider using Goth.
func ConfigureOIDCProvider() {
	clientID, clientSecret, callbackURL, endpoint := GetOIDCConfig()

	// Check if the OIDC configuration is set
	if clientID == "" || clientSecret == "" || callbackURL == "" || endpoint == "" {
		log.Println("OIDC configuration not set. Skipping OIDC configuration.")
		return
	}

	// Create a new OIDC provider
	provider, err := openidConnect.New(
		clientID,
		clientSecret,
		callbackURL,
		endpoint,
		"openid", "profile", "email",
	)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}

	// Set the name of the provider
	provider.SetName("oidc")

	// Add the provider to the list of Goth providers
	goth.UseProviders(provider)
}