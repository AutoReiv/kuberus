package rbac

import (
	"context"
	"errors"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RolesHandler handles role-related requests.
func RolesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if !isAdmin && !auth.HasPermission(username, "manage_roles") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage roles")
		}

		namespace := c.QueryParam("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch c.Request().Method {
		case http.MethodGet:
			return handleGetRoles(c, clientset, namespace)
		case http.MethodPost:
			return handleCreateRole(c, clientset, namespace)
		case http.MethodPut:
			return handleUpdateRole(c, clientset, namespace)
		case http.MethodDelete:
			return handleDeleteRole(c, clientset, namespace, c.QueryParam("name"))
		default:
			return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}
// IsRoleActive checks if a role is active by looking for any role bindings that reference it.
func IsRoleActive(clientset *kubernetes.Clientset, roleName, namespace string) (bool, error) {
	// Check RoleBindings in the namespace
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Error listing role bindings", zap.Error(err))
		return false, err
	}
	for _, rb := range roleBindings.Items {
		if rb.RoleRef.Name == roleName {
			return true, nil
		}
	}
	return false, nil
}

// RoleWithStatus represents a role with its active status.
type RoleWithStatus struct {
	rbacv1.Role
	Active bool `json:"active"`
}

// handleGetRoles handles listing roles in a specific namespace or across all namespaces.
func handleGetRoles(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	if namespace == "all" {
		return listAllNamespacesRoles(c, clientset)
	}
	return listNamespaceRoles(c, clientset, namespace)
}

// listNamespaceRoles lists roles in a specific namespace.
func listNamespaceRoles(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	roles, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Error listing roles", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var rolesWithStatus []RoleWithStatus
	for _, role := range roles.Items {
		active, err := IsRoleActive(clientset, role.Name, namespace)
		if err != nil {
			utils.Logger.Error("Error checking if role is active", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		rolesWithStatus = append(rolesWithStatus, RoleWithStatus{Role: role, Active: active})
	}

	utils.Logger.Info("Listed roles in namespace", zap.String("namespace", namespace))
	return c.JSON(http.StatusOK, rolesWithStatus)
}

// listAllNamespacesRoles lists roles across all namespaces.
func listAllNamespacesRoles(c echo.Context, clientset *kubernetes.Clientset) error {
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Error listing roles across all namespaces", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var rolesWithStatus []RoleWithStatus
	for _, role := range roles.Items {
		active, err := IsRoleActive(clientset, role.Name, role.Namespace)
		if err != nil {
			utils.Logger.Error("Error checking if role is active", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		rolesWithStatus = append(rolesWithStatus, RoleWithStatus{Role: role, Active: active})
	}

	utils.Logger.Info("Listed roles across all namespaces")
	return c.JSON(http.StatusOK, rolesWithStatus)
}

// handleCreateRole handles creating a new role in a specific namespace.
func handleCreateRole(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var role rbacv1.Role
	if err := c.Bind(&role); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind create role request")
	}

	if err := validateRole(&role); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid role", err, "Invalid role data")
	}

	createdRole, err := clientset.RbacV1().Roles(namespace).Create(context.TODO(), &role, metav1.CreateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to create role", err, "Failed to create role in Kubernetes")
	}

	utils.Logger.Info("Role created successfully", zap.String("roleName", role.Name), zap.String("namespace", namespace))
	utils.LogAuditEvent(c.Request(), "create", role.Name, namespace)
	return c.JSON(http.StatusOK, createdRole)
}

// handleUpdateRole handles updating an existing role in a specific namespace.
func handleUpdateRole(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var role rbacv1.Role
	if err := c.Bind(&role); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind update role request")
	}

	if err := validateRole(&role); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid role", err, "Invalid role data")
	}

	updatedRole, err := clientset.RbacV1().Roles(namespace).Update(context.TODO(), &role, metav1.UpdateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to update role", err, "Failed to update role in Kubernetes")
	}

	utils.Logger.Info("Role updated successfully", zap.String("roleName", role.Name), zap.String("namespace", namespace))
	utils.LogAuditEvent(c.Request(), "update", role.Name, namespace)
	return c.JSON(http.StatusOK, updatedRole)
}

// handleDeleteRole handles deleting a role in a specific namespace.
func handleDeleteRole(c echo.Context, clientset *kubernetes.Clientset, namespace, name string) error {
	if name == "" {
		utils.Logger.Warn("Role name is required")
		return echo.NewHTTPError(http.StatusBadRequest, "Role name is required")
	}

	err := clientset.RbacV1().Roles(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to delete role", err, "Failed to delete role in Kubernetes")
	}

	utils.Logger.Info("Role deleted successfully", zap.String("roleName", name), zap.String("namespace", namespace))
	utils.LogAuditEvent(c.Request(), "delete", name, namespace)
	return c.JSON(http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}

// RoleDetailsResponse represents the detailed information about a role.
type RoleDetailsResponse struct {
	Role         *rbacv1.Role         `json:"role"`
	RoleBindings []rbacv1.RoleBinding `json:"roleBindings"`
	Active       bool                 `json:"active"`
}

// RoleDetailsHandler handles fetching detailed information about a specific role.
func RoleDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if !isAdmin && !auth.HasPermission(username, "view_role_details") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view role details")
		}

		return getRoleDetails(c, clientset)
	}
}
// getRoleDetails fetches detailed information about a specific role

// getRoleDetails fetches detailed information about a specific role.
func getRoleDetails(c echo.Context, clientset *kubernetes.Clientset) error {
	roleName := c.QueryParam("roleName")
	namespace := c.QueryParam("namespace")
	if namespace == "" {
		namespace = "default"
	}

	role, err := clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Error fetching role details", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Error listing role bindings", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	associatedBindings := filterRoleBindings(roleBindings.Items, roleName)

	active, err := IsRoleActive(clientset, roleName, namespace)
	if err != nil {
		utils.Logger.Error("Error checking if role is active", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := RoleDetailsResponse{
		Role:         role,
		RoleBindings: associatedBindings,
		Active:       active,
	}

	utils.Logger.Info("Fetched role details", zap.String("roleName", roleName), zap.String("namespace", namespace))
	return c.JSON(http.StatusOK, response)
}

// filterRoleBindings filters role bindings associated with a specific role.
func filterRoleBindings(roleBindings []rbacv1.RoleBinding, roleName string) []rbacv1.RoleBinding {
	var associatedBindings []rbacv1.RoleBinding
	for _, rb := range roleBindings {
		if rb.RoleRef.Name == roleName {
			associatedBindings = append(associatedBindings, rb)
		}
	}
	return associatedBindings
}

// validateRole ensures that the role is valid.
func validateRole(role *rbacv1.Role) error {
	if role.Name == "" {
		return errors.New("role name is required")
	}
	if len(role.Rules) == 0 {
		return errors.New("at least one rule is required")
	}
	return nil
}
