package rbac

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceAccountsHandler handles requests related to service accounts.
func ServiceAccountsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch c.Request.Method {
		case http.MethodGet:
			listServiceAccounts(c, clientset, namespace)
		case http.MethodPost:
			createServiceAccount(c, clientset, namespace)
		case http.MethodDelete:
			deleteServiceAccount(c, clientset, namespace, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

func listServiceAccounts(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	serviceAccounts, err := clientset.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, serviceAccounts.Items)
}

func createServiceAccount(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	var serviceAccount corev1.ServiceAccount
	if err := c.ShouldBindJSON(&serviceAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body: " + err.Error()})
		return
	}

	createdServiceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &serviceAccount, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service account: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdServiceAccount)
}

func deleteServiceAccount(c *gin.Context, clientset *kubernetes.Clientset, namespace, name string) {
	err := clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service account: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}