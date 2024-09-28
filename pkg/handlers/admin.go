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

	// Check if there are any existing users
	var userCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error checking user count", err, "Failed to check user count")
	}

	// If there are existing users, check permissions
	if userCount > 0 {
		username, ok := c.Get("username").(string)
		if !ok || username == "" {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create an admin")
		}

		if !auth.HasPermission(username, "create_admin") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create an admin with permission: create_admin")
		}
	}

	var req CreateAdminRequest
	if err := c.Bind(&req); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request payload", err, "Failed to bind create admin request")
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
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error hashing password", err, "Failed to hash password")
	}

	// Check if an admin account already exists
	var adminExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&adminExists)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error checking admin existence", err, "Failed to check if admin exists")
	}

	if adminExists {
		utils.Logger.Warn("Admin account already exists")
		return echo.NewHTTPError(http.StatusConflict, "Admin account already exists")
	}

	// Store the admin account information
	_, err = db.DB.Exec("INSERT INTO users (username, password, source, is_admin) VALUES (?, ?, 'internal', true)", username, hashedPassword)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error creating admin account", err, "Failed to create admin account")
	}

	utils.Logger.Info("Admin account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "create_admin", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "Admin account created successfully"})
}
