package auth

import (
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

func ConfigureOIDCProvider() {
	provider, err := openidConnect.New(
		os.Getenv("OIDC_CLIENT_ID"),
		os.Getenv("OIDC_CLIENT_SECRET"),
		os.Getenv("OIDC_CALLBACK_URL"),
		os.Getenv("OIDC_ENDPOINT"),
		"openid", "profile", "email",
	)
	if err != nil {
		// Handle the error appropriately
		return
	}

	// Set the name of the provider
	provider.SetName("oidc")

	// Add the provider to the list of Goth providers
	goth.UseProviders(provider)
}
