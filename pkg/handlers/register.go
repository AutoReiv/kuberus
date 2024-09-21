package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RegisterRequest represents the request payload for user registration.
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

// RegisterResponse represents the response payload for a successful registration.
type RegisterResponse struct {
	Message string `json:"message"`
}

// RegisterHandler handles user registration.
func RegisterHandler(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		utils.Logger.Error("Invalid request payload", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
	}

	// Ensure Password and PasswordConfirm match
	if req.Password != req.PasswordConfirm {
		utils.Logger.Warn("Passwords do not match")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Passwords do not match"})
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		utils.Logger.Warn("Password does not meet strength requirements")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password does not meet strength requirements"})
	}

	// Create user using the user management handler logic
	if err := auth.CreateUser(username, password, "admin"); err != nil {
		utils.Logger.Error("Error creating user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user: " + err.Error()})
	}

	utils.Logger.Info("User account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "register_user", username, "N/A")
	return c.JSON(http.StatusOK, RegisterResponse{Message: "User account created successfully"})
}