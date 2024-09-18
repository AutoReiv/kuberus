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
		switch c.Request().Method {
		case http.MethodPost:
			return handleCreateUser(c)
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

func handleCreateUser(c echo.Context) error {
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

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		utils.Logger.Warn("Password does not meet strength requirements")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password does not meet strength requirements"})
	}

	// Create user
	if err := auth.CreateUser(username, password, "admin"); err != nil {
		utils.Logger.Error("Error creating user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user: " + err.Error()})
	}

	utils.Logger.Info("User account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "create_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account created successfully"})
}

func handleUpdateUser(c echo.Context) error {
	var req UpdateUserRequest
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

	// Update user
	if err := auth.UpdateUser(username, password); err != nil {
		utils.Logger.Error("Error updating user password", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user: " + err.Error()})
	}

	utils.Logger.Info("User account updated successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "update_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account updated successfully"})
}

func handleDeleteUser(c echo.Context) error {
	var req DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		utils.Logger.Error("Invalid request payload", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)

	// Delete user
	if err := auth.DeleteUser(username); err != nil {
		utils.Logger.Error("Error deleting user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting user: " + err.Error()})
	}

	utils.Logger.Info("User account deleted successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "delete_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account deleted successfully"})
}

func handleListUsers(c echo.Context) error {
	users, err := auth.GetAllUsers()
	if err != nil {
		utils.Logger.Error("Error retrieving users", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error retrieving users: " + err.Error()})
	}

	utils.Logger.Info("Retrieved user list successfully")
	return c.JSON(http.StatusOK, users)
}
