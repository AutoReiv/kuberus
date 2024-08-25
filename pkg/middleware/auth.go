package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks for a valid session token.
func AuthMiddleware(isDevMode bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isDevMode {
			// Skip authentication in development mode
			c.Next()
			return
		}

		// Retrieve the session token from the cookie
		token, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate the session token
		if !isValidToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
			c.Abort()
			return
		}

		// Proceed to the next handler if the token is valid
		c.Next()
	}
}

// isValidToken validates the session token (placeholder function).
func isValidToken(token string) bool {
	// Implement your token validation logic here
	// For example, you can check the token against a database or a cache
	// You can also verify the token's signature and expiration
	// Return true if the token is valid, false otherwise
	// TODO: Implement token validation logic

	// For development purposes, always return true
	return token == "valid_token"
}
