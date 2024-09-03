package middleware

import (
	"context"
	"log"
	"net/http"
	"rbac/pkg/auth"
	"strings"
	"time"
)

type contextKey string

const (
	usernameKey contextKey = "username"
)

// ApplyMiddlewares applies all the middlewares to the given handler.
func ApplyMiddlewares(handler http.Handler, isDevMode bool) http.Handler {
	if isDevMode {
		log.Println("Development mode: Applying middlewares")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingMiddleware(recoveryMiddleware(secureHeadersMiddleware(handler))).ServeHTTP(w, r)
	})
}

// loggingMiddleware logs the details of each HTTP request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

// recoveryMiddleware recovers from any panics and writes a 500 if there was one.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// secureHeadersMiddleware adds security-related headers to the response.
func secureHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware validates the JWT token and sets the user information in the request context.
func AuthMiddleware(next http.Handler, isDevMode bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevMode {
			// Debug statement to verify dev mode
			log.Println("Development mode: Bypassing authentication")
			// Skip authentication in development mode
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			log.Println("Token validation error:", err) // Debug statement
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set the username in the request context
		ctx := context.WithValue(r.Context(), usernameKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}