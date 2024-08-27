package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateAdminRequest represents the request payload for creating an admin account.
type CreateAdminRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

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

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not meet strength requirements"})
		return
	}
	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	// Acquire lock to synchronize access to shared data
	auth.Mu.Lock()
	defer auth.Mu.Unlock()

	// Check if an admin account already exists
	if auth.AdminExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Admin account already exists"})
		return
	}

	// Store the admin account information
	auth.Users[username] = hashedPassword
	auth.AdminExists = true

	c.JSON(http.StatusOK, gin.H{"message": "Admin account created successfully"})
}

// CreateUserRequest represents the request payload for creating a user account.
type CreateUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// CreateUserHandler handles the creation of a user account by the admin.
func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
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
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Acquire lock to synchronize access to shared data
	auth.Mu.Lock()
	defer auth.Mu.Unlock()

	// Check if the username already exists
	if _, exists := auth.Users[username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Store the user account information
	auth.Users[username] = hashedPassword

	c.JSON(http.StatusOK, gin.H{"message": "User account created successfully"})
}

type CreateOIDCConfigRequest struct {
	ClientID     string `json:"clientID" binding:"required"`
	ClientSecret string `json:"clientSecret" binding:"required"`
	CallbackURL  string `json:"callbackURL" binding:"required"`
	Endpoint     string `json:"endpoint" binding:"required"`
}

func CreateOIDCConfigHandler(c *gin.Context) {
	var req CreateOIDCConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store the OIDC configuration securely
	auth.StoreOIDCConfig(req.ClientID, req.ClientSecret, req.CallbackURL, req.Endpoint)

	c.JSON(http.StatusOK, gin.H{"message": "OIDC configuration created successfully"})
}
