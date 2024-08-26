package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RolesHandler handles role-related requests.
func RolesHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Query("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch c.Request.Method {
		case http.MethodGet:
			handleGetRoles(c, clientset, namespace)
		case http.MethodPost:
			handleCreateRole(c, clientset, namespace)
		case http.MethodPut:
			handleUpdateRole(c, clientset, namespace)
		case http.MethodDelete:
			handleDeleteRole(c, clientset, namespace, c.Query("name"))
		default:
			c.Status(http.StatusMethodNotAllowed)
		}
	}
}

// handleGetRoles handles listing roles in a specific namespace or across all namespaces.
func handleGetRoles(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	if namespace == "all" {
		listAllNamespacesRoles(c, clientset)
	} else {
		listNamespaceRoles(c, clientset, namespace)
	}
}

// listNamespaceRoles lists roles in a specific namespace.
func listNamespaceRoles(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
	roles, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// listAllNamespacesRoles lists roles across all namespaces.
func listAllNamespacesRoles(c *gin.Context, clientset *kubernetes.Clientset) {
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, roles)
}

// handleCreateRole handles creating a new role in a specific namespace.
func handleCreateRole(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
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

// handleUpdateRole handles updating an existing role in a specific namespace.
func handleUpdateRole(c *gin.Context, clientset *kubernetes.Clientset, namespace string) {
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

// handleDeleteRole handles deleting a role in a specific namespace.
func handleDeleteRole(c *gin.Context, clientset *kubernetes.Clientset, namespace, name string) {
	err := clientset.RbacV1().Roles(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

// RoleDetailsHandler handles fetching detailed information about a specific role.
func RoleDetailsHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		getRoleDetails(c, clientset)
	}
}

// getRoleDetails fetches detailed information about a specific role.
func getRoleDetails(c *gin.Context, clientset *kubernetes.Clientset) {
	roleName := c.Query("roleName")
	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = "default"
	}

	role, err := clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}

	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		handleError(c, err, http.StatusInternalServerError)
		return
	}

	associatedBindings := filterRoleBindings(roleBindings.Items, roleName)

	response := RoleDetailsResponse{
		Role:         role,
		RoleBindings: associatedBindings,
	}

	c.JSON(http.StatusOK, response)
}

// filterRoleBindings filters role bindings associated with a specific role.
func filterRoleBindings(roleBindings []rbacv1.RoleBinding, roleName string) []rbacv1.RoleBinding {
	var associatedBindings []rbacv1.RoleBinding
	for _, rb := range roleBindings {
		if rb.RoleRef.Name == roleName {
			associatedBindings = append(associatedBindings, rb)
		}
	}
	return associatedBindings
}

// RoleDetailsResponse represents the detailed information about a role.
type RoleDetailsResponse struct {
	Role         *rbacv1.Role         `json:"role"`
	RoleBindings []rbacv1.RoleBinding `json:"roleBindings"`
}

// handleError sends an error response with the specified status code.
func handleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}