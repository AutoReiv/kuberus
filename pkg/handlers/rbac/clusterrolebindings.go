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

// ClusterRoleBindingsHandler handles requests related to cluster role bindings.
func ClusterRoleBindingsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if (!isAdmin && !auth.HasPermission(username, "manage_clusterrolebindings")) {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage cluster role bindings")
		}

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
			return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

// handleListClusterRoleBindings lists all cluster role bindings.
func handleListClusterRoleBindings(c echo.Context, clientset *kubernetes.Clientset) error {
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing cluster role bindings", err, "Failed to list cluster role bindings")
	}
	utils.Logger.Info("Listed cluster role bindings")
	return c.JSON(http.StatusOK, clusterRoleBindings.Items)
}

// handleCreateClusterRoleBinding creates a new cluster role binding.
func handleCreateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := c.Bind(&clusterRoleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind create cluster role binding request")
	}

	if err := validateClusterRoleBinding(&clusterRoleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid cluster role binding", err, "Invalid cluster role binding data")
	}

	createdClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), &clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to create cluster role binding", err, "Failed to create cluster role binding in Kubernetes")
	}

	utils.Logger.Info("Cluster role binding created successfully", zap.String("clusterRoleBindingName", clusterRoleBinding.Name))
	utils.LogAuditEvent(c.Request(), "create", clusterRoleBinding.Name, "cluster-wide")
	return c.JSON(http.StatusOK, createdClusterRoleBinding)
}

// handleUpdateClusterRoleBinding updates an existing cluster role binding.
func handleUpdateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := c.Bind(&clusterRoleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind update cluster role binding request")
	}

	if err := validateClusterRoleBinding(&clusterRoleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid cluster role binding", err, "Invalid cluster role binding data")
	}

	updatedClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Update(context.TODO(), &clusterRoleBinding, metav1.UpdateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to update cluster role binding", err, "Failed to update cluster role binding in Kubernetes")
	}

	utils.Logger.Info("Cluster role binding updated successfully", zap.String("clusterRoleBindingName", clusterRoleBinding.Name))
	utils.LogAuditEvent(c.Request(), "update", clusterRoleBinding.Name, "cluster-wide")
	return c.JSON(http.StatusOK, updatedClusterRoleBinding)
}

// handleDeleteClusterRoleBinding deletes a cluster role binding by name.
func handleDeleteClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
    name := c.QueryParam("name")
    if name == "" {
        utils.Logger.Warn("Cluster role binding name is required")
        return echo.NewHTTPError(http.StatusBadRequest, "Cluster role binding name is required")
    }

    err := clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
    if err != nil {
        return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to delete cluster role binding", err, "Failed to delete cluster role binding in Kubernetes")
    }

    utils.Logger.Info("Cluster role binding deleted successfully", zap.String("clusterRoleBindingName", name))
    utils.LogAuditEvent(c.Request(), "delete", name, "cluster-wide")
    return c.JSON(http.StatusOK, map[string]string{"message": "Cluster role binding deleted successfully"})
}

// ClusterRoleBindingDetailsHandler handles fetching detailed information about a specific cluster role binding.
func ClusterRoleBindingDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if !isAdmin && !auth.HasPermission(username, "view_clusterrolebinding_details") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view cluster role binding details")
		}

		clusterRoleBindingName := c.QueryParam("name")
		if clusterRoleBindingName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Cluster role binding name is required")
		}

		clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBindingName, metav1.GetOptions{})
		if err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error fetching cluster role binding details", err, "Failed to fetch cluster role binding details")
		}

		utils.Logger.Info("Fetched cluster role binding details", zap.String("clusterRoleBindingName", clusterRoleBindingName))
		return c.JSON(http.StatusOK, clusterRoleBinding)
	}
}

// validateClusterRoleBinding ensures that the cluster role binding is valid.
func validateClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	if clusterRoleBinding.Name == "" {
		return errors.New("cluster role binding name is required")
	}
	if clusterRoleBinding.RoleRef.Name == "" {
		return errors.New("role reference name is required")
	}
	if len(clusterRoleBinding.Subjects) == 0 {
		return errors.New("at least one subject is required")
	}
	return nil
}
