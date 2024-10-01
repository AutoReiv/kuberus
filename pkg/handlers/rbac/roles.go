package rbac

import (
	"context"
	"net/http"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RolesHandler handles role-related requests.
func RolesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		if namespace == "" {
			namespace = "default"
		}

		handlers := map[string]func(echo.Context, *kubernetes.Clientset, string) error{
			http.MethodGet:    handleGetRoles,
			http.MethodPost:   handleCreateRole,
			http.MethodPut:    handleUpdateRole,
			http.MethodDelete: handleDeleteRole,
		}

		return utils.HandleHTTPMethod(c, clientset, namespace, handlers)
	}
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Error listing roles: "+err.Error())
	}

	var rolesWithStatus []RoleWithStatus
	for _, role := range roles.Items {
		active, err := IsRoleActive(clientset, role.Name, namespace)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error checking if role is active: "+err.Error())
		}
		rolesWithStatus = append(rolesWithStatus, RoleWithStatus{Role: role, Active: active})
	}

	return c.JSON(http.StatusOK, rolesWithStatus)
}

// listAllNamespacesRoles lists roles across all namespaces.
func listAllNamespacesRoles(c echo.Context, clientset *kubernetes.Clientset) error {
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error listing roles across all namespaces: "+err.Error())
	}

	var rolesWithStatus []RoleWithStatus
	for _, role := range roles.Items {
		active, err := IsRoleActive(clientset, role.Name, role.Namespace)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error checking if role is active: "+err.Error())
		}
		rolesWithStatus = append(rolesWithStatus, RoleWithStatus{Role: role, Active: active})
	}

	return c.JSON(http.StatusOK, rolesWithStatus)
}

// handleCreateRole handles creating a new role in a specific namespace.
func handleCreateRole(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var role rbacv1.Role
	if err := c.Bind(&role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request body: "+err.Error())
	}

	if err := utils.ValidateRole(&role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role: "+err.Error())
	}

	createdRole, err := clientset.RbacV1().Roles(namespace).Create(context.TODO(), &role, metav1.CreateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create role: "+err.Error())
	}

	return c.JSON(http.StatusOK, createdRole)
}

// handleUpdateRole handles updating an existing role in a specific namespace.
func handleUpdateRole(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var role rbacv1.Role
	if err := c.Bind(&role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request body: "+err.Error())
	}

	if err := utils.ValidateRole(&role); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role: "+err.Error())
	}

	updatedRole, err := clientset.RbacV1().Roles(namespace).Update(context.TODO(), &role, metav1.UpdateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update role: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedRole)
}

// handleDeleteRole handles deleting a role in a specific namespace.
func handleDeleteRole(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	name := c.QueryParam("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Role name is required")
	}

	err := clientset.RbacV1().Roles(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete role: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}

// RoleWithStatus represents a role with its active status.
type RoleWithStatus struct {
	rbacv1.Role
	Active bool `json:"active"`
}

// IsRoleActive checks if a role is active by looking for any role bindings that reference it.
func IsRoleActive(clientset *kubernetes.Clientset, roleName, namespace string) (bool, error) {
	// Check RoleBindings in the namespace
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, rb := range roleBindings.Items {
		if rb.RoleRef.Name == roleName {
			return true, nil
		}
	}
	return false, nil
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
		return getRoleDetails(c, clientset)
	}
}

// getRoleDetails fetches detailed information about a specific role.
func getRoleDetails(c echo.Context, clientset *kubernetes.Clientset) error {
	roleName := c.QueryParam("roleName")
	namespace := c.QueryParam("namespace")
	if namespace == "" {
		namespace = "default"
	}

	role, err := clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching role details: "+err.Error())
	}

	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error listing role bindings: "+err.Error())
	}

	associatedBindings := filterRoleBindings(roleBindings.Items, roleName)

	active, err := IsRoleActive(clientset, roleName, namespace)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking if role is active: "+err.Error())
	}

	response := RoleDetailsResponse{
		Role:         role,
		RoleBindings: associatedBindings,
		Active:       active,
	}

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
