package auth

import (
	"log"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

// ConfigureOIDCProvider sets up the OIDC provider using Goth.
func ConfigureOIDCProvider() {
	clientID := os.Getenv("OIDC_CLIENT_ID")
	clientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	callbackURL := os.Getenv("OIDC_CALLBACK_URL")
	endpoint := os.Getenv("OIDC_ENDPOINT")

	// Ensure all required environment variables are set
	if clientID == "" || clientSecret == "" || callbackURL == "" || endpoint == "" {
		log.Fatal("OIDC configuration error: missing required environment variables")
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
