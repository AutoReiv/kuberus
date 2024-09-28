package rbac

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

// ClusterRolesHandler handles requests related to cluster roles.
func ClusterRolesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if (!isAdmin && !auth.HasPermission(username, "manage_clusterroles")) {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage cluster roles")
		}

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
			return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}
// handleListClusterRoles lists all cluster roles.
func handleListClusterRoles(c echo.Context, clientset *kubernetes.Clientset) error {
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing cluster roles", err, "Failed to list cluster roles")
	}

	var clusterRolesWithStatus []ClusterRoleWithStatus
	for _, clusterRole := range clusterRoles.Items {
		active, err := IsClusterRoleActive(clientset, clusterRole.Name)
		if err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error checking if cluster role is active", err, "Failed to check if cluster role is active")
		}
		clusterRolesWithStatus = append(clusterRolesWithStatus, ClusterRoleWithStatus{ClusterRole: clusterRole, Active: active})
	}

	utils.Logger.Info("Listed cluster roles")
	return c.JSON(http.StatusOK, clusterRolesWithStatus)
}

// handleCreateClusterRole creates a new cluster role.
func handleCreateClusterRole(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRole rbacv1.ClusterRole
	if err := c.Bind(&clusterRole); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind create cluster role request")
	}

	createdClusterRole, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(), &clusterRole, metav1.CreateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to create cluster role", err, "Failed to create cluster role in Kubernetes")
	}

	utils.Logger.Info("Cluster role created successfully", zap.String("clusterRoleName", clusterRole.Name))
	utils.LogAuditEvent(c.Request(), "create", clusterRole.Name, "cluster-wide")
	return c.JSON(http.StatusOK, createdClusterRole)
}

// handleUpdateClusterRole updates an existing cluster role.
func handleUpdateClusterRole(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRole rbacv1.ClusterRole
	if err := c.Bind(&clusterRole); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind update cluster role request")
	}

	updatedClusterRole, err := clientset.RbacV1().ClusterRoles().Update(context.TODO(), &clusterRole, metav1.UpdateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to update cluster role", err, "Failed to update cluster role in Kubernetes")
	}

	utils.Logger.Info("Cluster role updated successfully", zap.String("clusterRoleName", clusterRole.Name))
	utils.LogAuditEvent(c.Request(), "update", clusterRole.Name, "cluster-wide")
	return c.JSON(http.StatusOK, updatedClusterRole)
}

// handleDeleteClusterRole deletes a cluster role by name.
func handleDeleteClusterRole(c echo.Context, clientset *kubernetes.Clientset) error {
    name := c.QueryParam("name")
    if name == "" {
        utils.Logger.Warn("Cluster role name is required")
        return echo.NewHTTPError(http.StatusBadRequest, "Cluster role name is required")
    }

    err := clientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
    if err != nil {
        return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to delete cluster role", err, "Failed to delete cluster role in Kubernetes")
    }

    utils.Logger.Info("Cluster role deleted successfully", zap.String("clusterRoleName", name))
    utils.LogAuditEvent(c.Request(), "delete", name, "cluster-wide")
    return c.JSON(http.StatusOK, map[string]string{"message": "Cluster role deleted successfully"})
}

// ClusterRoleDetailsHandler handles fetching detailed information about a specific cluster role.
func ClusterRoleDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		if !auth.HasPermission(username, "view_clusterrole_details") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view cluster role details")
		}

		return handleGetClusterRoleDetails(c, clientset)
	}
}

// handleGetClusterRoleDetails fetches detailed information about a specific cluster role.
func handleGetClusterRoleDetails(c echo.Context, clientset *kubernetes.Clientset) error {
	clusterRoleName := c.QueryParam("clusterRoleName")
	if clusterRoleName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Cluster role name is required")
	}

	clusterRole, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), clusterRoleName, metav1.GetOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error fetching cluster role details", err, "Failed to fetch cluster role details")
	}

	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing cluster role bindings", err, "Failed to list cluster role bindings")
	}

	associatedBindings := filterClusterRoleBindings(clusterRoleBindings.Items, clusterRoleName)

	active, err := IsClusterRoleActive(clientset, clusterRoleName)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error checking if cluster role is active", err, "Failed to check if cluster role is active")
	}

	response := ClusterRoleDetailsResponse{
		ClusterRole:         clusterRole,
		ClusterRoleBindings: associatedBindings,
		Active:              active,
	}

	utils.Logger.Info("Fetched cluster role details", zap.String("clusterRoleName", clusterRoleName))
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

// ClusterRoleWithStatus represents a cluster role with its active status.
type ClusterRoleWithStatus struct {
	rbacv1.ClusterRole
	Active bool `json:"active"`
}

// ClusterRoleDetailsResponse represents the detailed information about a cluster role.
type ClusterRoleDetailsResponse struct {
	ClusterRole         *rbacv1.ClusterRole         `json:"clusterRole"`
	ClusterRoleBindings []rbacv1.ClusterRoleBinding `json:"clusterRoleBindings"`
	Active              bool                        `json:"active"`
}

// IsClusterRoleActive checks if a cluster role is active by looking for any cluster role bindings that reference it.
func IsClusterRoleActive(clientset *kubernetes.Clientset, clusterRoleName string) (bool, error) {
	// Check ClusterRoleBindings
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Error listing cluster role bindings", zap.Error(err))
		return false, err
	}
	for _, crb := range clusterRoleBindings.Items {
		if crb.RoleRef.Name == clusterRoleName {
			return true, nil
		}
	}
	return false, nil
}
