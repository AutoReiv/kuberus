package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

// CreateUserRequest represents the request payload for creating a user.
type CreateUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Role            string `json:"role" binding:"required"` // New field for role
}

// UpdateUserRequest represents the request payload for updating a user's password.
type UpdateUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

// DeleteUserRequest represents the request payload for deleting a user.
type DeleteUserRequest struct {
	Username string `json:"username" binding:"required"`
}

// UserManagementHandler handles user management-related requests.
func UserManagementHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if !isAdmin && !auth.HasPermission(username, "manage_users") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage users")
		}

		switch c.Request().Method {
		case http.MethodPost:
			return HandleCreateUser(c, clientset)
		case http.MethodPut:
			return handleUpdateUser(c)
		case http.MethodDelete:
			return handleDeleteUser(c)
		case http.MethodGet:
			return handleListUsers(c)
		default:
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	}
}

// HandleCreateUser creates a new user and assigns a specified role.
func HandleCreateUser(c echo.Context, clientset *kubernetes.Clientset) error {
	requestingUsername := c.Get("username").(string)
	if !auth.HasPermission(requestingUsername, "create_user") {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create users")
	}

	var req CreateUserRequest
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
	role := utils.SanitizeInput(req.Role)

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		utils.Logger.Warn("Password does not meet strength requirements")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password does not meet strength requirements"})
	}

	// Create user
	if err := auth.CreateUser(username, password, "internal"); err != nil {
		utils.Logger.Error("Error creating user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user: " + err.Error()})
	}

	// Assign the specified role to the new user
	if err := auth.AssignRoleToUser(username, role); err != nil {
		utils.Logger.Error("Error assigning role to user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error assigning role to user: " + err.Error()})
	}

	utils.Logger.Info("User account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "create_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account created successfully"})
}

// handleUpdateUser updates an existing user's details.
func handleUpdateUser(c echo.Context) error {
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

func handleDeleteUser(c echo.Context) error {
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

func handleListUsers(c echo.Context) error {
	users, err := auth.GetAllUsers()
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error retrieving users", err, "Failed to retrieve users")
	}

	utils.Logger.Info("Retrieved user list successfully")
	return c.JSON(http.StatusOK, users)
}
