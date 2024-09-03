package middleware

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWT secret key
var jwtKey []byte

// Initialize the JWT secret key
func init() {
	jwtKey = generateRandomKey(32) // 32 bytes for HS256
}

// generateRandomKey generates a secure random key of the specified length
func generateRandomKey(length int) []byte {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
	}
	return key
}

// Claims defines the structure of the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// ApplyMiddlewares applies all the middlewares to the given handler.
func ApplyMiddlewares(handler http.Handler, isDevMode bool) http.Handler {
	if isDevMode {
		log.Println("Development mode: Applying middlewares")
	}
	return loggingMiddleware(recoveryMiddleware(secureHeadersMiddleware(handler)))
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

// AuthMiddleware checks for a valid session token.
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer " prefix
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims, err := validateJWT(tokenStr)
		if err != nil {
			log.Println("Token validation error:", err) // Debug statement
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store the username in the context
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateJWT validates a JWT token and returns the claims if valid
func validateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
