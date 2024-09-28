package handlers

import (
	"context"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		username := c.Get("username").(string)
		if !auth.HasPermission(username, "manage_users") {
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
// HandleCreateUser creates a new user and assigns a read-only role.
func HandleCreateUser(c echo.Context, clientset *kubernetes.Clientset) error {
	username := c.Get("username").(string)
	if !auth.HasPermission(username, "create_user") {
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
	username = utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

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

	// Ensure the read-only role exists
	if err := EnsureReadOnlyRole(clientset, "default"); err != nil {
		utils.Logger.Error("Error ensuring read-only role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error ensuring read-only role: " + err.Error()})
	}

	// Assign read-only role to the new user
	if err := assignReadOnlyRole(clientset, username); err != nil {
		utils.Logger.Error("Error assigning read-only role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error assigning read-only role: " + err.Error()})
	}

	utils.Logger.Info("User account created successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "create_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account created successfully"})
}

// EnsureReadOnlyRole ensures that the read-only role exists in the specified namespace.
func EnsureReadOnlyRole(clientset *kubernetes.Clientset, namespace string) error {
	roleName := "read-only"

	// Check if the role already exists
	_, err := clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err == nil {
		// Role already exists, no need to create it
		return nil
	}

	// Define the read-only role
	readOnlyRole := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services", "deployments"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}

	// Create the read-only role
	_, err = clientset.RbacV1().Roles(namespace).Create(context.TODO(), readOnlyRole, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to create read-only role", zap.Error(err))
		return err
	}

	utils.Logger.Info("Read-only role created successfully", zap.String("roleName", roleName), zap.String("namespace", namespace))
	return nil
}

// assignReadOnlyRole assigns the read-only role to a user.
func assignReadOnlyRole(clientset *kubernetes.Clientset, username string) error {
	// Create a RoleBinding to assign the read-only role to the user
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      username + "-read-only",
			Namespace: "default",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: "read-only",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	_, err := clientset.RbacV1().RoleBindings("default").Create(context.TODO(), roleBinding, metav1.CreateOptions{})
	return err
}

func handleUpdateUser(c echo.Context) error {
	username := c.Get("username").(string)
	if !auth.HasPermission(username, "update_user") {
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
	username = utils.SanitizeInput(req.Username)
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
	username := c.Get("username").(string)
	if !auth.HasPermission(username, "delete_user") {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to delete users")
	}

	var req DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request payload", err, "Failed to bind delete user request")
	}

	// Sanitize user input
	username = utils.SanitizeInput(req.Username)

	// Delete user
	if err := auth.DeleteUser(username); err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error deleting user", err, "Failed to delete user account")
	}

	utils.Logger.Info("User account deleted successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "delete_user", username, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "User account deleted successfully"})

}

// handleListUsers handles the retrieval of user accounts
func handleListUsers(c echo.Context) error {
	username := c.Get("username").(string)
	if !auth.HasPermission(username, "list_users") {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to list users")
	}

	users, err := auth.GetAllUsers()
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error retrieving users", err, "Failed to retrieve users")
	}

	utils.Logger.Info("Retrieved user list successfully")
	return c.JSON(http.StatusOK, users)
}