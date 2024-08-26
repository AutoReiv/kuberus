package utils

import (
	"crypto/rand"
	"encoding/base64"
	"html"
	"regexp"
	"unicode"
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

// IsStrongPassword checks if a password meets strength requirements.
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case isSpecialChar(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// isSpecialChar checks if a character is a special character.
func isSpecialChar(char rune) bool {
	specialChars := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	return specialChars.MatchString(string(char))
}
