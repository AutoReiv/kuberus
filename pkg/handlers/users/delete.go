
package users

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// DeleteUserRequest represents the request payload for deleting a user.
type DeleteUserRequest struct {
	Username string `json:"username" binding:"required"`
}

// HandleDeleteUser deletes a user.
func HandleDeleteUser(c echo.Context) error {
	requestingUsername := c.Get("username").(string)
	if !auth.HasPermission(requestingUsername, "delete_user") {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to delete users")
	}

	var req DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request payload", err, "Failed to bind delete user request")
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)

	// Delete user
	if err := auth.DeleteUser(username); err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error deleting user", err, "Failed to delete user account")
	}

	utils.Logger.Info("User account deleted successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "delete_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account deleted successfully"})
}
