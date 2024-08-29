package middleware

import (
	"fmt"
	"net/http"
	"rbac/pkg/auth"
	"strings"

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

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract the token from the "Bearer " prefix
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			fmt.Println("Token validation error:", err) // Debug statement
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store the username in the context
		c.Set("username", claims.Username)
		c.Next()
	}
}