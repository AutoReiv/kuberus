package rbac

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceAccountsHandler handles requests related to service accounts.
func ServiceAccountsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, serviceAccounts.Items)
}

// handleCreateServiceAccount creates a new service account in a specific namespace.
func handleCreateServiceAccount(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var serviceAccount corev1.ServiceAccount
	if err := c.Bind(&serviceAccount); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to decode request body: " + err.Error()})
	}

	createdServiceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &serviceAccount, metav1.CreateOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create service account: " + err.Error()})
	}

	return c.JSON(http.StatusOK, createdServiceAccount)
}

// handleDeleteServiceAccount deletes a service account in a specific namespace.
func handleDeleteServiceAccount(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Service account name is required"})
	}

	err := clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete service account: " + err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
