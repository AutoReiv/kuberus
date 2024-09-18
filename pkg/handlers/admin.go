package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/db"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CreateAdminRequest represents the request payload for creating an admin account.
type CreateAdminRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// CreateAdminHandler handles the creation of an admin account.
func CreateAdminHandler(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	var req CreateAdminRequest
	if err := c.Bind(&req); err != nil {
		utils.Logger.Error("Invalid request payload", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
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

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		utils.Logger.Error("Error hashing password", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Error hashing password: "+err.Error())
	}

	// Check if an admin account already exists
	var adminExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&adminExists)
	if err != nil {
		utils.Logger.Error("Error checking admin existence", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking admin existence: "+err.Error())
	}

	if adminExists {
		utils.Logger.Warn("Admin account already exists")
		return echo.NewHTTPError(http.StatusConflict, "Admin account already exists")
	}

	// Store the admin account information
	_, err = db.DB.Exec("INSERT INTO users (username, password, source) VALUES (?, ?, 'internal')", username, hashedPassword)
	if err != nil {
		utils.Logger.Error("Error creating admin account", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating admin account: "+err.Error())
	}

	utils.Logger.Info("Admin account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "create_admin", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "Admin account created successfully"})
}