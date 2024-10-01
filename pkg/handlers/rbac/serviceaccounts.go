package rbac

import (
	"context"
	"net/http"
	"rbac/pkg/utils"

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

		handlers := map[string]func(echo.Context, *kubernetes.Clientset, string) error{
			http.MethodGet:    handleListServiceAccounts,
			http.MethodPost:   handleCreateServiceAccount,
			http.MethodDelete: handleDeleteServiceAccount,
		}

		return utils.HandleHTTPMethod(c, clientset, namespace, handlers)
	}
}

// handleListServiceAccounts lists all service accounts in a specific namespace.
func handleListServiceAccounts(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	listFunc := func(namespace string, opts metav1.ListOptions) (interface{}, error) {
		return clientset.CoreV1().ServiceAccounts(namespace).List(context.TODO(), opts)
	}
	return utils.ListResources(c, clientset, namespace, listFunc)
}

// handleCreateServiceAccount creates a new service account in a specific namespace.
func handleCreateServiceAccount(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var serviceAccount corev1.ServiceAccount
	createFunc := func(namespace string, obj interface{}, opts metav1.CreateOptions) (interface{}, error) {
		return clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), obj.(*corev1.ServiceAccount), opts)
	}
	return utils.CreateResource(c, clientset, namespace, &serviceAccount, createFunc)
}

// handleDeleteServiceAccount deletes a service account in a specific namespace.
func handleDeleteServiceAccount(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	name := c.QueryParam("name")
	deleteFunc := func(namespace, name string, opts metav1.DeleteOptions) error {
		return clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, opts)
	}
	return utils.DeleteResource(c, clientset, namespace, name, deleteFunc)
}