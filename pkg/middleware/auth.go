package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks for a valid session token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate the session token (this is a placeholder, implement your own validation)
		if !isValidToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isValidToken validates the session token (placeholder function).
func isValidToken(token string) bool {
	// Implement your token validation logic here
	return token == "valid_token"
}
