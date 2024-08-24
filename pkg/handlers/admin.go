package handlers

import (
	"net/http"
	"rbac/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateAdminRequest represents the request payload for creating an admin account.
type CreateAdminRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

var adminExists bool

// CreateAdminHandler handles the creation of an admin account.
func CreateAdminHandler(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Acquire lock to synchronize access to shared data
	mu.Lock()
	defer mu.Unlock()

	// Check if an admin account already exists
	if adminExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Admin account already exists"})
		return
	}

	// Store the admin account information
	users[username] = hashedPassword
	adminExists = true

	c.JSON(http.StatusOK, gin.H{"message": "Admin account created successfully"})
}
