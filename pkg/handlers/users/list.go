package users

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
)

// HandleListUsers lists all users.
func HandleListUsers(c echo.Context) error {
	users, err := auth.GetAllUsers()
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error retrieving users", err, "Failed to retrieve users")
	}

	utils.Logger.Info("Retrieved user list successfully")
	return c.JSON(http.StatusOK, users)
}
