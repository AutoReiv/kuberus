package rbac

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ClusterRoleBindingsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet:
			listClusterRoleBindings(c, clientset)
		case http.MethodPost:
			createClusterRoleBindings(c, clientset)
		case http.MethodDelete:
			deleteClusterRoleBinding(c, clientset, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

func listClusterRoleBindings(c *gin.Context, clientset *kubernetes.Clientset) {
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clusterRoleBindings.Items)
}

func createClusterRoleBindings(c *gin.Context, clientset *kubernetes.Clientset) {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := c.ShouldBindJSON(&clusterRoleBinding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body: " + err.Error()})
		return
	}

	createdClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), &clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cluster role binding: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdClusterRoleBinding)
}

func deleteClusterRoleBinding(c *gin.Context, clientset *kubernetes.Clientset, name string) {
	err := clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cluster role binding: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
