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

// NamespacesHandler handles requests related to namespaces.
func NamespacesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		handlers := map[string]func(echo.Context, *kubernetes.Clientset, string) error{
			http.MethodGet:    handleListNamespaces,
			http.MethodPost:   handleCreateNamespace,
			http.MethodDelete: handleDeleteNamespace,
		}

		return utils.HandleHTTPMethod(c, clientset, "", handlers)
	}
}

// handleListNamespaces lists all namespaces.
func handleListNamespaces(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	return utils.ListResources(c, clientset, "", func(namespace string, opts metav1.ListOptions) (interface{}, error) {
		return clientset.CoreV1().Namespaces().List(context.TODO(), opts)
	})
}

// handleCreateNamespace creates a new namespace.
func handleCreateNamespace(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	var namespace corev1.Namespace
	return utils.CreateResource(c, clientset, "", &namespace, func(namespace string, obj interface{}, opts metav1.CreateOptions) (interface{}, error) {
		return clientset.CoreV1().Namespaces().Create(context.TODO(), obj.(*corev1.Namespace), opts)
	})
}

// handleDeleteNamespace deletes a namespace by name.
func handleDeleteNamespace(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	name := c.QueryParam("name")
	return utils.DeleteResource(c, clientset, "", name, func(namespace, name string, opts metav1.DeleteOptions) error {
		return clientset.CoreV1().Namespaces().Delete(context.TODO(), name, opts)
	})
}