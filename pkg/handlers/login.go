package handlers

import (
	"net/http"
	"rbac/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the request payload for user login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response payload for a successful login.
type LoginResponse struct {
	Message string `json:"message"`
}

// LoginHandler handles user login.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

	// Implement rate limiting to prevent brute-force attacks
	if !rateLimit(c.ClientIP()) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		return
	}

	// Acquire lock to synchronize access to shared data
	mu.Lock()
	hashedPassword, ok := users[username]
	mu.Unlock()

	// Authenticate user
	if !ok || !utils.CheckPasswordHash(password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate and store secure session token
	sessionToken := utils.GenerateSessionToken()
	session := &Session{
		Username: username,
		Token:    sessionToken,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}
	mu.Lock()
	sessions[sessionToken] = session
	mu.Unlock()

	// Set a secure cookie for session management
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie.Name, cookie.Value, int(cookie.Expires.Unix()), cookie.Path, "", cookie.Secure, cookie.HttpOnly)

	c.JSON(http.StatusOK, LoginResponse{Message: "Login successful"})
}
