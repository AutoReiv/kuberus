package auth

import (
	"crypto/rand"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey []byte

// Initialize the JWT secret key
func init() {
	JwtKey = generateRandomKey(32) // 32 bytes for HS256
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
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token for a given username
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	// Debug statement
	log.Printf("Generated JWT: %s", tokenString)
	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims if valid
func ValidateJWT(tokenStr string) (*Claims, error) {
	// Debug statement
	log.Printf("Validating JWT: %s", tokenStr)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		// Debug statement
		log.Printf("Validation error: %v", err)
		return nil, err
	}

	if !token.Valid {
		// Debug statement
		log.Println("Invalid token signature")
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
