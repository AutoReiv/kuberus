package rbac

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"rbac/pkg/utils"
)

// NamespacesHandler handles requests related to namespaces.
func NamespacesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		switch c.Request().Method {
		case http.MethodGet:
			return handleListNamespaces(c, clientset)
		case http.MethodPost:
			return handleCreateNamespace(c, clientset)
		case http.MethodDelete:
			return handleDeleteNamespace(c, clientset)
		default:
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	}
}

// handleListNamespaces lists all namespaces.
func handleListNamespaces(c echo.Context, clientset *kubernetes.Clientset) error {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Error listing namespaces", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	utils.Logger.Info("Listed namespaces")
	return c.JSON(http.StatusOK, namespaces.Items)
}

// handleCreateNamespace creates a new namespace.
func handleCreateNamespace(c echo.Context, clientset *kubernetes.Clientset) error {
	var namespace corev1.Namespace
	if err := c.Bind(&namespace); err != nil {
		utils.Logger.Error("Failed to decode request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	createdNamespace, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to create namespace", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create namespace: " + err.Error()})
	}

	utils.Logger.Info("Namespace created successfully", zap.String("namespaceName", namespace.Name))
	utils.LogAuditEvent(c.Request(), "create", namespace.Name, "N/A")
	return c.JSON(http.StatusOK, createdNamespace)
}

// handleDeleteNamespace deletes a namespace by name.
func handleDeleteNamespace(c echo.Context, clientset *kubernetes.Clientset) error {
	name := c.QueryParam("name")
	if name == "" {
		utils.Logger.Warn("Namespace name is required")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Namespace name is required"})
	}

	err := clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to delete namespace", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete namespace: " + err.Error()})
	}

	utils.Logger.Info("Namespace deleted successfully", zap.String("namespaceName", name))
	utils.LogAuditEvent(c.Request(), "delete", name, "N/A")
	return c.NoContent(http.StatusNoContent)
}