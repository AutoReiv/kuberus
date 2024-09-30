package auth

import (
	"net/http"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// CheckPermission checks if a user has a specific permission.
func CheckPermission(username, permission string) bool {
	if IsAdmin(username) {
		return true
	}
	return HasPermission(username, permission)
}

// LogAndRespondPermissionDenied logs and responds with a permission denied error.
func LogAndRespondPermissionDenied(c echo.Context, username, permission string) error {
	utils.Logger.Warn("Permission denied", zap.String("username", username), zap.String("permission", permission))
	return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to access this resource")
}
