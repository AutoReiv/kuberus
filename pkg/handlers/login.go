package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiter = rate.NewLimiter(1, 5) // Allow 1 request per second with a burst of 5
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
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
				return
			}

			// Authenticate user
			if !auth.AuthenticateUser(username, password) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			// Generate and store secure session token
			sessionToken := auth.GenerateSessionToken()
			session := &auth.Session{
				Username: username,
				Token:    sessionToken,
				ExpireAt: time.Now().Add(24 * time.Hour),
			}
			auth.StoreSession(session)

			// Set a secure cookie for session management
			expiration := time.Now().Add(24 * time.Hour)
			c.SetCookie("session_token", sessionToken, int(expiration.Unix()), "/", "", true, true)

			c.JSON(http.StatusOK, LoginResponse{Message: "Login successful"})
		}
// authenticateUser authenticates the user credentials.
func authenticateUser(username, password string) bool {
	// Acquire lock to synchronize access to shared data
	auth.Mu.Lock()
	defer auth.Mu.Unlock()

	hashedPassword, ok := auth.Users[username]
	if !ok {
		return false
	}

	return auth.CheckPasswordHash(password, hashedPassword)
}
