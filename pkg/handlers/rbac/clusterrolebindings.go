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

// ClusterRoleBindingsHandler handles requests related to cluster role bindings.
func ClusterRoleBindingsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		switch c.Request().Method {
		case http.MethodGet:
			return handleListClusterRoleBindings(c, clientset)
		case http.MethodPost:
			return handleCreateClusterRoleBinding(c, clientset)
		case http.MethodPut:
			return handleUpdateClusterRoleBinding(c, clientset)
		case http.MethodDelete:
			return handleDeleteClusterRoleBinding(c, clientset)
		default:
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	}
}

// handleListClusterRoleBindings lists all cluster role bindings.
func handleListClusterRoleBindings(c echo.Context, clientset *kubernetes.Clientset) error {
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, clusterRoleBindings.Items)
}

// handleCreateClusterRoleBinding creates a new cluster role binding.
func handleCreateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := c.Bind(&clusterRoleBinding); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	createdClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), &clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create cluster role binding: " + err.Error()})
	}

	utils.LogAuditEvent(c.Request(), "create", clusterRoleBinding.Name, "cluster-wide")
	return c.JSON(http.StatusOK, createdClusterRoleBinding)
}

// handleUpdateClusterRoleBinding updates an existing cluster role binding.
func handleUpdateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := c.Bind(&clusterRoleBinding); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	updatedClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Update(context.TODO(), &clusterRoleBinding, metav1.UpdateOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update cluster role binding: " + err.Error()})
	}

	utils.LogAuditEvent(c.Request(), "update", clusterRoleBinding.Name, "cluster-wide")
	return c.JSON(http.StatusOK, updatedClusterRoleBinding)
}

// handleDeleteClusterRoleBinding deletes a cluster role binding by name.
func handleDeleteClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cluster role binding name is required"})
	}

	err := clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete cluster role binding: " + err.Error()})
	}

	utils.LogAuditEvent(c.Request(), "delete", name, "cluster-wide")
	return c.NoContent(http.StatusNoContent)
}

// ClusterRoleBindingDetailsHandler handles fetching detailed information about a specific cluster role binding.
func ClusterRoleBindingDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		clusterRoleBindingName := c.QueryParam("name")
		if clusterRoleBindingName == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cluster role binding name is required"})
		}

		clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBindingName, metav1.GetOptions{})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, clusterRoleBinding)
	}
}
