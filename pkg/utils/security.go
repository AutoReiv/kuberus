package utils

import (
	"crypto/rand"
	"encoding/base64"
	"html"
	"golang.org/x/crypto/bcrypt"
)

// SanitizeInput sanitizes the input string to prevent XSS attacks.
func SanitizeInput(input string) string {
	return html.EscapeString(input)
}

// GenerateSessionToken generates a secure random session token.
func GenerateSessionToken() string {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	return base64.URLEncoding.EncodeToString(tokenBytes)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
