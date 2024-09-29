package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/db"
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
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request payload", err, "Failed to bind register request")
	}

	// Ensure Password and PasswordConfirm match
	if req.Password != req.PasswordConfirm {
		utils.Logger.Warn("Passwords do not match")
		return echo.NewHTTPError(http.StatusBadRequest, "Passwords do not match")
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		utils.Logger.Warn("Password does not meet strength requirements")
		return echo.NewHTTPError(http.StatusBadRequest, "Password does not meet strength requirements")
	}

	// Check if there are any existing users
	var userCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error checking user count", err, "Failed to check user count")
	}

	// Create user
	if userCount == 0 {
		// First user, create as admin
		if err := auth.CreateUser(username, password, "internal"); err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error creating user", err, "Failed to create user account")
		}
		// Assign admin role to the first user
		if err := auth.AssignRoleToUser(username, "admin"); err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error assigning admin role", err, "Failed to assign admin role to user")
		}
	} else {
		// Create as regular user
		if err := auth.CreateUser(username, password, "internal"); err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error creating user", err, "Failed to create user account")
		}
		// Assign viewer role to subsequent users
		if err := auth.AssignRoleToUser(username, "viewer"); err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error assigning viewer role", err, "Failed to assign viewer role to user")
		}
	}

	utils.Logger.Info("User account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "register_user", username, "N/A")
	return c.JSON(http.StatusOK, RegisterResponse{Message: "User account created successfully"})
}