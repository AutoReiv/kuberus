package rbac

import (
	"context"
	"net/http"

	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceAccountsHandler handles requests related to service accounts.
func ServiceAccountsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		if !auth.HasPermission(username, "manage_serviceaccounts") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage service accounts")
		}

		namespace := c.QueryParam("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch c.Request().Method {
		case http.MethodGet:
			return handleListServiceAccounts(c, clientset, namespace)
		case http.MethodPost:
			return handleCreateServiceAccount(c, clientset, namespace)
		case http.MethodDelete:
			return handleDeleteServiceAccount(c, clientset, namespace)
		default:
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	}
}

// handleListServiceAccounts lists all service accounts in a specific namespace.
func handleListServiceAccounts(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	serviceAccounts, err := clientset.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing service accounts", err, "Failed to list service accounts")
	}
	utils.Logger.Info("Listed service accounts", zap.String("namespace", namespace))
	return c.JSON(http.StatusOK, serviceAccounts.Items)
}

// handleCreateServiceAccount creates a new service account in a specific namespace.
func handleCreateServiceAccount(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var serviceAccount corev1.ServiceAccount
	if err := c.Bind(&serviceAccount); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind create service account request")
	}

	createdServiceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &serviceAccount, metav1.CreateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to create service account", err, "Failed to create service account in Kubernetes")
	}

	utils.Logger.Info("Service account created successfully", zap.String("serviceAccountName", serviceAccount.Name), zap.String("namespace", namespace))
	utils.LogAuditEvent(c.Request(), "create", serviceAccount.Name, namespace)
	return c.JSON(http.StatusOK, createdServiceAccount)
}

// handleDeleteServiceAccount deletes a service account in a specific namespace.
func handleDeleteServiceAccount(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	name := c.QueryParam("name")
	if name == "" {
		utils.Logger.Warn("Service account name is required")
		return echo.NewHTTPError(http.StatusBadRequest, "Service account name is required")
	}

	err := clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to delete service account", err, "Failed to delete service account in Kubernetes")
	}

	utils.Logger.Info("Service account deleted successfully", zap.String("serviceAccountName", name), zap.String("namespace", namespace))
	utils.LogAuditEvent(c.Request(), "delete", name, namespace)
	return c.NoContent(http.StatusNoContent)
}
