package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// HandleHTTPMethod handles different HTTP methods for a given handler function.
func HandleHTTPMethod(c echo.Context, clientset *kubernetes.Clientset, namespace string, handlers map[string]func(echo.Context, *kubernetes.Clientset, string) error) error {
	if handler, exists := handlers[c.Request().Method]; exists {
		return handler(c, clientset, namespace)
	}
	return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
}

// ListResources lists resources in a specific namespace.
func ListResources(c echo.Context, clientset *kubernetes.Clientset, namespace string, listFunc func(string, metav1.ListOptions) (interface{}, error)) error {
	resources, err := listFunc(namespace, metav1.ListOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error listing resources: "+err.Error())
	}
	return c.JSON(http.StatusOK, resources)
}

// CreateResource creates a new resource in a specific namespace.
func CreateResource(c echo.Context, clientset *kubernetes.Clientset, namespace string, resource interface{}, createFunc func(string, interface{}, metav1.CreateOptions) (interface{}, error)) error {
	if err := c.Bind(resource); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request body: "+err.Error())
	}

	createdResource, err := createFunc(namespace, resource, metav1.CreateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create resource: "+err.Error())
	}

	return c.JSON(http.StatusOK, createdResource)
}

// UpdateResource updates an existing resource in a specific namespace.
func UpdateResource(c echo.Context, clientset *kubernetes.Clientset, namespace string, resource interface{}, updateFunc func(string, interface{}, metav1.UpdateOptions) (interface{}, error)) error {
	if err := c.Bind(resource); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request body: "+err.Error())
	}

	updatedResource, err := updateFunc(namespace, resource, metav1.UpdateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update resource: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedResource)
}

// DeleteResource deletes a resource by name in a specific namespace.
func DeleteResource(c echo.Context, clientset *kubernetes.Clientset, namespace, name string, deleteFunc func(string, string, metav1.DeleteOptions) error) error {
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Resource name is required")
	}

	err := deleteFunc(namespace, name, metav1.DeleteOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete resource: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Resource deleted successfully"})
}
