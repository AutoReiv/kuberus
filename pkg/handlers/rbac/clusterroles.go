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

// ClusterRolesHandler handles requests related to cluster roles.
func ClusterRolesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		switch c.Request().Method {
		case http.MethodGet:
			return handleListClusterRoles(c, clientset)
		case http.MethodPost:
			return handleCreateClusterRole(c, clientset)
		case http.MethodPut:
			return handleUpdateClusterRole(c, clientset)
		case http.MethodDelete:
			return handleDeleteClusterRole(c, clientset)
		default:
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	}
}

// IsClusterRoleActive checks if a cluster role is active by looking for any cluster role bindings that reference it.
func IsClusterRoleActive(clientset *kubernetes.Clientset, clusterRoleName string) (bool, error) {
	// Check ClusterRoleBindings
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, crb := range clusterRoleBindings.Items {
		if crb.RoleRef.Name == clusterRoleName {
			return true, nil
		}
	}
	return false, nil
}

// ClusterRoleWithStatus represents a cluster role with its active status.
type ClusterRoleWithStatus struct {
	rbacv1.ClusterRole
	Active bool `json:"active"`
}

// handleListClusterRoles lists all cluster roles.
func handleListClusterRoles(c echo.Context, clientset *kubernetes.Clientset) error {
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var clusterRolesWithStatus []ClusterRoleWithStatus
	for _, clusterRole := range clusterRoles.Items {
		active, err := IsClusterRoleActive(clientset, clusterRole.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		clusterRolesWithStatus = append(clusterRolesWithStatus, ClusterRoleWithStatus{ClusterRole: clusterRole, Active: active})
	}

	return c.JSON(http.StatusOK, clusterRolesWithStatus)
}

// handleCreateClusterRole creates a new cluster role.
func handleCreateClusterRole(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRole rbacv1.ClusterRole
	if err := c.Bind(&clusterRole); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	createdClusterRole, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(), &clusterRole, metav1.CreateOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create cluster role: " + err.Error()})
	}

	utils.LogAuditEvent(c.Request(), "create", clusterRole.Name, "cluster-wide")
	return c.JSON(http.StatusOK, createdClusterRole)
}

// handleUpdateClusterRole updates an existing cluster role.
func handleUpdateClusterRole(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRole rbacv1.ClusterRole
	if err := c.Bind(&clusterRole); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	updatedClusterRole, err := clientset.RbacV1().ClusterRoles().Update(context.TODO(), &clusterRole, metav1.UpdateOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update cluster role: " + err.Error()})
	}

	utils.LogAuditEvent(c.Request(), "update", clusterRole.Name, "cluster-wide")
	return c.JSON(http.StatusOK, updatedClusterRole)
}

// handleDeleteClusterRole deletes a cluster role by name.
func handleDeleteClusterRole(c echo.Context, clientset *kubernetes.Clientset) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cluster role name is required"})
	}

	err := clientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete cluster role: " + err.Error()})
	}

	utils.LogAuditEvent(c.Request(), "delete", name, "cluster-wide")
	return c.NoContent(http.StatusNoContent)
}

// ClusterRoleDetailsHandler handles fetching detailed information about a specific cluster role.
func ClusterRoleDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleGetClusterRoleDetails(c, clientset)
	}
}

// ClusterRoleDetailsResponse represents the detailed information about a cluster role.
type ClusterRoleDetailsResponse struct {
	ClusterRole         *rbacv1.ClusterRole         `json:"clusterRole"`
	ClusterRoleBindings []rbacv1.ClusterRoleBinding `json:"clusterRoleBindings"`
	Active              bool                        `json:"active"`
}

// handleGetClusterRoleDetails fetches detailed information about a specific cluster role.
func handleGetClusterRoleDetails(c echo.Context, clientset *kubernetes.Clientset) error {
	clusterRoleName := c.QueryParam("clusterRoleName")
	if clusterRoleName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cluster role name is required"})
	}

	clusterRole, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), clusterRoleName, metav1.GetOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	associatedBindings := filterClusterRoleBindings(clusterRoleBindings.Items, clusterRoleName)

	active, err := IsClusterRoleActive(clientset, clusterRoleName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := ClusterRoleDetailsResponse{
		ClusterRole:         clusterRole,
		ClusterRoleBindings: associatedBindings,
		Active:              active,
	}

	return c.JSON(http.StatusOK, response)
}

// filterClusterRoleBindings filters cluster role bindings associated with a specific cluster role.
func filterClusterRoleBindings(clusterRoleBindings []rbacv1.ClusterRoleBinding, clusterRoleName string) []rbacv1.ClusterRoleBinding {
	var associatedBindings []rbacv1.ClusterRoleBinding
	for _, crb := range clusterRoleBindings {
		if crb.RoleRef.Name == clusterRoleName {
			associatedBindings = append(associatedBindings, crb)
		}
	}
	return associatedBindings
}
