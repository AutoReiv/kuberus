package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	IssuerURL    string
	RedirectURL  string
}

var verifier *oidc.IDTokenVerifier
var oauth2Config oauth2.Config

func InitOIDC(config OIDCConfig) error {
	provider, err := oidc.NewProvider(context.Background(), config.IssuerURL)
	if err != nil {
		return fmt.Errorf("failed to get provider: %v", err)
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	oauth2Config = oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  config.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return nil
}

func OIDCMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := verifier.Verify(context.Background(), token)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		// Extract custom claims if needed
		var claims struct {
			Email string `json:"email"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusUnauthorized)
			return
		}

		// Add claims to context if needed
		type contextKey string
		ctx := context.WithValue(r.Context(), contextKey("email"), claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
