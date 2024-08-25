package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespacesHandler handles listing namespaces
func NamespacesHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet:
			listNamespaces(c, clientset)
		case http.MethodPost:
			createNamespace(c, clientset)
		case http.MethodDelete:
			deleteNamespace(c, clientset, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

// listNamespaces lists all namespaces
func listNamespaces(c *gin.Context, clientset *kubernetes.Clientset) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, namespaces.Items)
}

// createNamespace creates a new namespace
func createNamespace(c *gin.Context, clientset *kubernetes.Clientset) {
	var namespace corev1.Namespace
	if err := c.ShouldBindJSON(&namespace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body: " + err.Error()})
		return
	}

	createdNamespace, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create namespace: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdNamespace)
}

// deleteNamespace deletes a namespace
func deleteNamespace(c *gin.Context, clientset *kubernetes.Clientset, name string) {
	if err := clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete namespace: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
