package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
)

// AuthMiddleware verifies the OIDC token and extracts user information
func AuthMiddleware(verifier *oidc.IDTokenVerifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		idToken, err := verifier.Verify(context.Background(), token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var claims struct {
			Email string `json:"email"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Add user information to the request context
		ctx := context.WithValue(r.Context(), "user", claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
