package rbac

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RoleBindingsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch c.Request.Method {
		case http.MethodGet:
			listRoleBindings(c, clientset, namespace)
		case http.MethodPost:
			createRoleBindings(c, clientset, namespace)
		case http.MethodDelete:
			deleteRoleBinding(c, clientset, namespace, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

func listRoleBindings(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roleBindings.Items)
}

func createRoleBindings(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	var roleBinding rbacv1.RoleBinding
	if err := c.ShouldBindJSON(&roleBinding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body: " + err.Error()})
		return
	}

	createdRoleBinding, err := clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), &roleBinding, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role binding: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdRoleBinding)
}

func deleteRoleBinding(c *gin.Context, clientset *kubernetes.Clientset, namespace, name string) {
	err := clientset.RbacV1().RoleBindings(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role binding: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
