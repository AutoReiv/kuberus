package handlers

import (
	"net/http"
	"rbac/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RegisterRequest represents the request payload for user registration.
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// RegisterResponse represents the response payload for a successful registration.
type RegisterResponse struct {
	Message string `json:"message"`
}

// RegisterHandler handles user registration.
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not meet strength requirements"})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Acquire lock to synchronize access to shared data
	mu.Lock()
	defer mu.Unlock()

	// Check if the username already exists
	if _, exists := users[username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Store the user account information
	users[username] = hashedPassword

	c.JSON(http.StatusOK, RegisterResponse{Message: "User registered successfully"})
}
