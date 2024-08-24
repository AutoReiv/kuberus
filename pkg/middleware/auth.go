package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("session_token")
	if err != nil || cookie != "authenticated" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}
