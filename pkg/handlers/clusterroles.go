package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ClusterRolesHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet:
			listClusterRoles(c, clientset)
		case http.MethodPost:
			createClusterRole(c, clientset)
		case http.MethodDelete:
			deleteClusterRole(c, clientset, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

func listClusterRoles(c *gin.Context, clientset *kubernetes.Clientset) {
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clusterRoles.Items)
}

func createClusterRole(c *gin.Context, clientset *kubernetes.Clientset) {
	var clusterRole rbacv1.ClusterRole
	if err := c.ShouldBindJSON(&clusterRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body: " + err.Error()})
		return
	}

	createdClusterRole, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(), &clusterRole, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cluster role: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdClusterRole)
}

func deleteClusterRole(c *gin.Context, clientset *kubernetes.Clientset, name string) {
	err := clientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cluster role: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
