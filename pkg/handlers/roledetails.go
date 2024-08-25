package handlers

import (
	"context"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/gin-gonic/gin"
)

// roleDetailsHandler handles fetching detailed information about a specific role
func RoleDetailsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		getRoleDetails(c, clientset)
	}
}

// getRoleDetails fetches detailed information about a specific role
func getRoleDetails(c *gin.Context, clientset *kubernetes.Clientset) {
	// Get the role name and namespace from the query parameters
	roleName := c.Query("roleName")
	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// Fetch the Role details
	role, err := clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the associated RoleBindings
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter RoleBindings that are associated with the role
	var associatedBindings []rbacv1.RoleBinding
	for _, rb := range roleBindings.Items {
		if rb.RoleRef.Name == roleName {
			associatedBindings = append(associatedBindings, rb)
		}
	}

	// Create a response structure
	type RoleDetailsResponse struct {
		Role         *rbacv1.Role         `json:"role"`
		RoleBindings []rbacv1.RoleBinding `json:"roleBindings"`
		// UsageStatistics can be added here if needed
	}

	response := RoleDetailsResponse{
		Role:         role,
		RoleBindings: associatedBindings,
	}

	// Return the detailed role information as JSON
	c.JSON(http.StatusOK, response)
}
