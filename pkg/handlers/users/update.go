package users

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// UpdateUserRequest represents the request payload for updating a user's password.
type UpdateUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

// HandleUpdateUser updates an existing user's details.
func HandleUpdateUser(c echo.Context) error {
	requestingUsername := c.Get("username").(string)
	if !auth.HasPermission(requestingUsername, "update_user") {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to update users")
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request payload", err, "Failed to bind update user request")
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

	// Update user
	if err := auth.UpdateUser(username, password); err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error updating user password", err, "Failed to update user password")
	}

	utils.Logger.Info("User account updated successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "update_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account updated successfully"})
}
