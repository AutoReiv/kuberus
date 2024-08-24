package handlers

import (
	"net/http"
	"rbac/pkg/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

// Mock user data store
var (
	users       = map[string]string{}
	adminExists = false
	mu          sync.Mutex
)

// RegisterRequest represents the registration request payload.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse represents the registration response payload.
type RegisterResponse struct {
	Message string `json:"message"`
}

// RegisterHandler handles user registration.
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if adminExists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registration is closed"})
		return
	}

	if _, exists := users[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	users[req.Username] = hashedPassword
	adminExists = true

	c.JSON(http.StatusOK, RegisterResponse{Message: "User registered successfully"})
}
