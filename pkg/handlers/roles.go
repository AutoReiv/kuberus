package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RolesHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch c.Request.Method {
		case http.MethodGet:
			if namespace == "all" {
				listRolesAllNamespaces(c, clientset)
			} else {
				listRoles(c, clientset, namespace)
			}
		case http.MethodPost:
			createRole(c, clientset, namespace)
		case http.MethodPut:
			editRole(c, clientset, namespace)
		case http.MethodDelete:
			deleteRole(c, clientset, namespace, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

func listRoles(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	roles, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, roles)
}

func listRolesAllNamespaces(c *gin.Context, clientset *kubernetes.Clientset) {
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, roles)
}

func createRole(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	var role rbacv1.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		handleError(c, err, http.StatusBadRequest)
		return
	}

	createdRole, err := clientset.RbacV1().Roles(namespace).Create(context.TODO(), &role, metav1.CreateOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, createdRole)
}

func editRole(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	var role rbacv1.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		handleError(c, err, http.StatusBadRequest)
		return
	}

	updatedRole, err := clientset.RbacV1().Roles(namespace).Update(context.TODO(), &role, metav1.UpdateOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedRole)
}

func deleteRole(c *gin.Context, clientset *kubernetes.Clientset, namespace, name string) {
	err := clientset.RbacV1().Roles(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func handleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
