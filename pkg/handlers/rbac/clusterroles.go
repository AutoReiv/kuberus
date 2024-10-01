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
		handlers := map[string]func(echo.Context, *kubernetes.Clientset, string) error{
			http.MethodGet:    handleListClusterRoles,
			http.MethodPost:   handleCreateClusterRole,
			http.MethodPut:    handleUpdateClusterRole,
			http.MethodDelete: handleDeleteClusterRole,
		}

		return utils.HandleHTTPMethod(c, clientset, "", handlers)
	}
}

// handleListClusterRoles lists all cluster roles.
func handleListClusterRoles(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	return utils.ListResources(c, clientset, "", func(namespace string, opts metav1.ListOptions) (interface{}, error) {
		return clientset.RbacV1().ClusterRoles().List(context.TODO(), opts)
	})
}

// handleCreateClusterRole creates a new cluster role.
func handleCreateClusterRole(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	var clusterRole rbacv1.ClusterRole
	return utils.CreateResource(c, clientset, "", &clusterRole, func(namespace string, obj interface{}, opts metav1.CreateOptions) (interface{}, error) {
		return clientset.RbacV1().ClusterRoles().Create(context.TODO(), obj.(*rbacv1.ClusterRole), opts)
	})
}

// handleUpdateClusterRole updates an existing cluster role.
func handleUpdateClusterRole(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	var clusterRole rbacv1.ClusterRole
	return utils.UpdateResource(c, clientset, "", &clusterRole, func(namespace string, obj interface{}, opts metav1.UpdateOptions) (interface{}, error) {
		return clientset.RbacV1().ClusterRoles().Update(context.TODO(), obj.(*rbacv1.ClusterRole), opts)
	})
}

// handleDeleteClusterRole deletes a cluster role by name.
func handleDeleteClusterRole(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	name := c.QueryParam("name")
	return utils.DeleteResource(c, clientset, "", name, func(namespace, name string, opts metav1.DeleteOptions) error {
		return clientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, opts)
	})
}

// ClusterRoleDetailsHandler handles fetching detailed information about a specific cluster role.
func ClusterRoleDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching cluster role details: "+err.Error())
	}

	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error listing cluster role bindings: "+err.Error())
	}

	associatedBindings := filterClusterRoleBindings(clusterRoleBindings.Items, clusterRoleName)

	active, err := IsClusterRoleActive(clientset, clusterRoleName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking if cluster role is active: "+err.Error())
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
		return false, err
	}
	for _, crb := range clusterRoleBindings.Items {
		if crb.RoleRef.Name == clusterRoleName {
			return true, nil
		}
	}
	return false, nil
}
