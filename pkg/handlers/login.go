package handlers

import (
	"net/http"
	"rbac/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response payload.
type LoginResponse struct {
	Message string `json:"message"`
}

// LoginHandler handles the login page.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Authenticate user
	hashedPassword, ok := users[req.Username]
	if !ok || !utils.CheckPasswordHash(req.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Set a cookie for session management
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "session_token", Value: "authenticated", Expires: expiration}
	c.SetCookie(cookie.Name, cookie.Value, int(cookie.Expires.Unix()), cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

	c.JSON(http.StatusOK, LoginResponse{Message: "Login successful"})
}
