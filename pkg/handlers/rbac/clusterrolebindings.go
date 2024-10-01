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

// ClusterRoleBindingsHandler handles requests related to cluster role bindings.
func ClusterRoleBindingsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		handlers := map[string]func(echo.Context, *kubernetes.Clientset, string) error{
			http.MethodGet:    handleListClusterRoleBindings,
			http.MethodPost:   handleCreateClusterRoleBinding,
			http.MethodPut:    handleUpdateClusterRoleBinding,
			http.MethodDelete: handleDeleteClusterRoleBinding,
		}

		return utils.HandleHTTPMethod(c, clientset, "", handlers)
	}
}

// handleListClusterRoleBindings lists all cluster role bindings.
func handleListClusterRoleBindings(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	return utils.ListResources(c, clientset, "", func(namespace string, opts metav1.ListOptions) (interface{}, error) {
		return clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), opts)
	})
}

// handleCreateClusterRoleBinding creates a new cluster role binding.
func handleCreateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	return utils.CreateResource(c, clientset, "", &clusterRoleBinding, func(namespace string, obj interface{}, opts metav1.CreateOptions) (interface{}, error) {
		return clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), obj.(*rbacv1.ClusterRoleBinding), opts)
	})
}

// handleUpdateClusterRoleBinding updates an existing cluster role binding.
func handleUpdateClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	return utils.UpdateResource(c, clientset, "", &clusterRoleBinding, func(namespace string, obj interface{}, opts metav1.UpdateOptions) (interface{}, error) {
		return clientset.RbacV1().ClusterRoleBindings().Update(context.TODO(), obj.(*rbacv1.ClusterRoleBinding), opts)
	})
}

// handleDeleteClusterRoleBinding deletes a cluster role binding by name.
func handleDeleteClusterRoleBinding(c echo.Context, clientset *kubernetes.Clientset, _ string) error {
	name := c.QueryParam("name")
	return utils.DeleteResource(c, clientset, "", name, func(namespace, name string, opts metav1.DeleteOptions) error {
		return clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, opts)
	})
}

// ClusterRoleBindingDetailsHandler handles fetching detailed information about a specific cluster role binding.
func ClusterRoleBindingDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		clusterRoleBindingName := c.QueryParam("name")
		if clusterRoleBindingName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Cluster role binding name is required")
		}

		clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBindingName, metav1.GetOptions{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching cluster role binding details: "+err.Error())
		}

		return c.JSON(http.StatusOK, clusterRoleBinding)
	}
}