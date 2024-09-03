package middleware

import (
	"context"
	"log"
	"net/http"
	"rbac/pkg/auth"
	"strings"
	"time"
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
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add username to context
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
