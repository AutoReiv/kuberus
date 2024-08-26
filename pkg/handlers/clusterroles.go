package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ClusterRolesHandler handles requests related to cluster roles.
func ClusterRolesHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet:
			handleListClusterRoles(c, clientset)
		case http.MethodPost:
			handleCreateClusterRole(c, clientset)
		case http.MethodDelete:
			handleDeleteClusterRole(c, clientset, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

// handleListClusterRoles lists all cluster roles.
func handleListClusterRoles(c *gin.Context, clientset *kubernetes.Clientset) {
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clusterRoles.Items)
}

// handleCreateClusterRole creates a new cluster role.
func handleCreateClusterRole(c *gin.Context, clientset *kubernetes.Clientset) {
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

// handleDeleteClusterRole deletes a cluster role by name.
func handleDeleteClusterRole(c *gin.Context, clientset *kubernetes.Clientset, name string) {
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster role name is required"})
		return
	}

	err := clientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cluster role: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ClusterRoleDetailsHandler handles fetching detailed information about a specific cluster role.
func ClusterRoleDetailsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		handleGetClusterRoleDetails(c, clientset)
	}
}

// handleGetClusterRoleDetails fetches detailed information about a specific cluster role.
func handleGetClusterRoleDetails(c *gin.Context, clientset *kubernetes.Clientset) {
	clusterRoleName := c.Query("clusterRoleName")
	if clusterRoleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster role name is required"})
		return
	}

	clusterRole, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), clusterRoleName, metav1.GetOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	associatedBindings := filterClusterRoleBindings(clusterRoleBindings.Items, clusterRoleName)

	response := ClusterRoleDetailsResponse{
		ClusterRole:         clusterRole,
		ClusterRoleBindings: associatedBindings,
	}

	c.JSON(http.StatusOK, response)
}

// filterClusterRoleBindings filters cluster role bindings associated with a specific cluster role.
func filterClusterRoleBindings(clusterRoleBindings []rbacv1.ClusterRoleBinding, clusterRoleName string) []rbacv1.ClusterRoleBinding {
	var associatedBindings []rbacv1.ClusterRoleBinding
	for _, crb := range clusterRoleBindings {
		if crb.RoleRef.Name == clusterRoleName {
			associatedBindings = append(associatedBindings, crb)
		}
	}
	return associatedBindings
}

// ClusterRoleDetailsResponse represents the detailed information about a cluster role.
type ClusterRoleDetailsResponse struct {
	ClusterRole         *rbacv1.ClusterRole         `json:"clusterRole"`
	ClusterRoleBindings []rbacv1.ClusterRoleBinding `json:"clusterRoleBindings"`
}
