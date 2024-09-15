package rbac

import (
	"context"
	"errors"
	"net/http"

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
		utils.Logger.Error("Error listing cluster role bindings", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	utils.Logger.Info("Listed cluster role bindings")
	return c.JSON(http.StatusOK, clusterRoleBindings.Items)
}

// handleCreateClusterRoleBinding creates a new cluster role binding.
func handleCreateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := c.Bind(&clusterRoleBinding); err != nil {
		utils.Logger.Error("Failed to decode request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	if err := validateClusterRoleBinding(&clusterRoleBinding); err != nil {
		utils.Logger.Error("Invalid cluster role binding", zap.Error(err))
		return c.JSON(http.StatusBadRequest, "Invalid cluster role binding: "+err.Error())
	}

	createdClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), &clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to create cluster role binding", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create cluster role binding: " + err.Error()})
	}

	utils.Logger.Info("Cluster role binding created successfully", zap.String("clusterRoleBindingName", clusterRoleBinding.Name))
	utils.LogAuditEvent(c.Request(), "create", clusterRoleBinding.Name, "cluster-wide")
	return c.JSON(http.StatusOK, createdClusterRoleBinding)
}

// handleUpdateClusterRoleBinding updates an existing cluster role binding.
func handleUpdateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := c.Bind(&clusterRoleBinding); err != nil {
		utils.Logger.Error("Failed to decode request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	if err := validateClusterRoleBinding(&clusterRoleBinding); err != nil {
		utils.Logger.Error("Invalid cluster role binding", zap.Error(err))
		return c.JSON(http.StatusBadRequest, "Invalid cluster role binding: "+err.Error())
	}

	updatedClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Update(context.TODO(), &clusterRoleBinding, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to update cluster role binding", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update cluster role binding: " + err.Error()})
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cluster role binding name is required"})
	}

	err := clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to delete cluster role binding", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete cluster role binding: " + err.Error()})
	}

	utils.Logger.Info("Cluster role binding deleted successfully", zap.String("clusterRoleBindingName", name))
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
			utils.Logger.Error("Error fetching cluster role binding details", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
