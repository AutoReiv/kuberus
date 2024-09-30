package users

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

// HandleCreateUser creates a new user and assigns a specified role.
func HandleCreateUser(c echo.Context, clientset *kubernetes.Clientset) error {
	requestingUsername := c.Get("username").(string)
	if !auth.HasPermission(requestingUsername, "create_user") {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to create users")
	}

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request payload", err, "Failed to bind create user request")
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
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error creating user", err, "Failed to create user account")
	}

	// Assign the specified role to the new user
	if err := auth.AssignRoleToUser(username, role); err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error assigning role to user", err, "Failed to assign role to user")
	}

	utils.Logger.Info("User account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "create_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account created successfully"})
}
