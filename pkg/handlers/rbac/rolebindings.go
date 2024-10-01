package rbac

import (
	"context"
	"net/http"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RoleBindingsHandler handles role binding-related requests.
func RoleBindingsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		if namespace == "" {
			namespace = "default"
		}

		handlers := map[string]func(echo.Context, *kubernetes.Clientset, string) error{
			http.MethodGet:    handleListRoleBindings,
			http.MethodPost:   handleCreateRoleBinding,
			http.MethodPut:    handleUpdateRoleBinding,
			http.MethodDelete: handleDeleteRoleBinding,
		}

		return utils.HandleHTTPMethod(c, clientset, namespace, handlers)
	}
}

// handleListRoleBindings lists all role bindings in a specific namespace.
func handleListRoleBindings(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	return utils.ListResources(c, clientset, namespace, func(namespace string, opts metav1.ListOptions) (interface{}, error) {
		return clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), opts)
	})
}

// handleCreateRoleBinding creates a new role binding in a specific namespace.
func handleCreateRoleBinding(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var roleBinding rbacv1.RoleBinding
	return utils.CreateResource(c, clientset, namespace, &roleBinding, func(namespace string, obj interface{}, opts metav1.CreateOptions) (interface{}, error) {
		return clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), obj.(*rbacv1.RoleBinding), opts)
	})
}

// handleUpdateRoleBinding updates an existing role binding in a specific namespace.
func handleUpdateRoleBinding(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var roleBinding rbacv1.RoleBinding
	return utils.UpdateResource(c, clientset, namespace, &roleBinding, func(namespace string, obj interface{}, opts metav1.UpdateOptions) (interface{}, error) {
		return clientset.RbacV1().RoleBindings(namespace).Update(context.TODO(), obj.(*rbacv1.RoleBinding), opts)
	})
}

// handleDeleteRoleBinding deletes a role binding in a specific namespace.
func handleDeleteRoleBinding(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	name := c.QueryParam("name")
	return utils.DeleteResource(c, clientset, namespace, name, func(namespace, name string, opts metav1.DeleteOptions) error {
		return clientset.RbacV1().RoleBindings(namespace).Delete(context.TODO(), name, opts)
	})
}

// RoleBindingDetailsHandler handles fetching detailed information about a specific role binding.
func RoleBindingDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		roleBindingName := c.QueryParam("name")
		namespace := c.QueryParam("namespace")
		if namespace == "" {
			namespace = "default"
		}

		roleBinding, err := clientset.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBindingName, metav1.GetOptions{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching role binding details: "+err.Error())
		}

		return c.JSON(http.StatusOK, roleBinding)
	}
}
