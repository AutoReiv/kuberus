package rbac

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"rbac/pkg/auth"
	"rbac/pkg/utils"
)

// NamespacesHandler handles requests related to namespaces.
func NamespacesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if !isAdmin && !auth.HasPermission(username, "manage_namespaces") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage namespaces")
		}

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
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing namespaces", err, "Failed to list namespaces")
	}
	utils.Logger.Info("Listed namespaces")
	return c.JSON(http.StatusOK, namespaces.Items)
}

// handleCreateNamespace creates a new namespace.
func handleCreateNamespace(c echo.Context, clientset *kubernetes.Clientset) error {
	var namespace corev1.Namespace
	if err := c.Bind(&namespace); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind create namespace request")
	}

	createdNamespace, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to create namespace", err, "Failed to create namespace in Kubernetes")
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
        return echo.NewHTTPError(http.StatusBadRequest, "Namespace name is required")
    }

    err := clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
    if err != nil {
        return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to delete namespace", err, "Failed to delete namespace in Kubernetes")
    }

    utils.Logger.Info("Namespace deleted successfully", zap.String("namespaceName", name))
    utils.LogAuditEvent(c.Request(), "delete", name, "N/A")
    return c.JSON(http.StatusOK, map[string]string{"message": "Namespace deleted successfully"})
}