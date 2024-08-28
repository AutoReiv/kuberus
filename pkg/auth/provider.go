package auth

import (
	"crypto/rand"
	"log"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

// GenerateRandomKey generates a random key of the given length.
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

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

	// Generate a random key for the session store
	key, err := GenerateRandomKey(32) // 32 bytes = 256 bits
	if err != nil {
		log.Fatalf("Failed to generate random key: %v", err)
	}

	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true in production

	store := sessions.NewCookieStore(key)
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should be enabled
	store.Options.Secure = isProd // Set to true in production

	gothic.Store = store
}
