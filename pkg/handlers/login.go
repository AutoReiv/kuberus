package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoginRequest represents the request payload for user login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response payload for a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles user login.
func LoginHandler(c echo.Context) error {
    var req LoginRequest
    if err := c.Bind(&req); err != nil {
        utils.Logger.Error("Invalid request payload", zap.Error(err))
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
    }

    // Sanitize user input
    username := utils.SanitizeInput(req.Username)
    password := utils.SanitizeInput(req.Password)

    // Authenticate user
    if !auth.AuthenticateUser(username, password) {
        utils.Logger.Warn("Login failed", zap.String("username", username))
        utils.LogAuditEvent(c.Request(), "login_failed", username, "N/A")
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
    }

    // Check if the user is an admin
    isAdmin := auth.IsAdmin(username)

    // Generate JWT token
    token, err := auth.GenerateJWT(username, isAdmin)
    if err != nil {
        utils.Logger.Error("Failed to generate token", zap.Error(err))
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token: " + err.Error()})
    }

    utils.Logger.Info("Login successful", zap.String("username", username))
    utils.LogAuditEvent(c.Request(), "login_success", username, "N/A")
    return c.JSON(http.StatusOK, LoginResponse{Token: token})
}
