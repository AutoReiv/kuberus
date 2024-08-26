package middleware

import (
	"fmt"
	"net/http"
	"rbac/pkg/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks for a valid session token.
func AuthMiddleware(isDevMode bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isDevMode {
			// Debug statement to verify dev mode
			fmt.Println("Development mode: Bypassing authentication")
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

// isValidToken validates the session token.
func isValidToken(token string) bool {
	session, exists := auth.GetSession(token)
	return exists && session != nil
}
