package rbac

import (
	"context"
	"errors"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RoleBindingsHandler handles role binding-related requests.
func RoleBindingsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		isAdmin, ok := c.Get("isAdmin").(bool)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
		}

		if !isAdmin && !auth.HasPermission(username, "manage_rolebindings") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to manage role bindings")
		}

		namespace := c.QueryParam("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch c.Request().Method {
		case http.MethodGet:
			return handleListRoleBindings(c, clientset, namespace)
		case http.MethodPost:
			return handleCreateRoleBinding(c, clientset, namespace)
		case http.MethodPut:
			return handleUpdateRoleBinding(c, clientset, namespace)
		case http.MethodDelete:
			return handleDeleteRoleBinding(c, clientset, namespace, c.QueryParam("name"))
		default:
			return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}
// handleListRoleBindings lists all role bindings in a specific namespace.
func handleListRoleBindings(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing role bindings", err, "Failed to list role bindings")
	}
	utils.Logger.Info("Listed role bindings", zap.String("namespace", namespace))
	return c.JSON(http.StatusOK, roleBindings.Items)
}

// handleCreateRoleBinding creates a new role binding in a specific namespace.
func handleCreateRoleBinding(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var roleBinding rbacv1.RoleBinding
	if err := c.Bind(&roleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind create role binding request")
	}

	if err := validateRoleBinding(&roleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid role binding", err, "Invalid role binding data")
	}

	createdRoleBinding, err := clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), &roleBinding, metav1.CreateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to create role binding", err, "Failed to create role binding in Kubernetes")
	}

	utils.Logger.Info("Role binding created successfully", zap.String("roleBindingName", roleBinding.Name), zap.String("namespace", namespace))
	utils.LogAuditEvent(c.Request(), "create", roleBinding.Name, namespace)
	return c.JSON(http.StatusOK, createdRoleBinding)
}

// handleUpdateRoleBinding updates an existing role binding in a specific namespace.
func handleUpdateRoleBinding(c echo.Context, clientset *kubernetes.Clientset, namespace string) error {
	var roleBinding rbacv1.RoleBinding
	if err := c.Bind(&roleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Failed to decode request body", err, "Failed to bind update role binding request")
	}

	if err := validateRoleBinding(&roleBinding); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid role binding", err, "Invalid role binding data")
	}

	updatedRoleBinding, err := clientset.RbacV1().RoleBindings(namespace).Update(context.TODO(), &roleBinding, metav1.UpdateOptions{})
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to update role binding", err, "Failed to update role binding in Kubernetes")
	}

	utils.Logger.Info("Role binding updated successfully", zap.String("roleBindingName", roleBinding.Name), zap.String("namespace", namespace))
	utils.LogAuditEvent(c.Request(), "update", roleBinding.Name, namespace)
	return c.JSON(http.StatusOK, updatedRoleBinding)
}

// handleDeleteRoleBinding deletes a role binding in a specific namespace.
func handleDeleteRoleBinding(c echo.Context, clientset *kubernetes.Clientset, namespace, name string) error {
    if name == "" {
        utils.Logger.Warn("Role binding name is required")
        return echo.NewHTTPError(http.StatusBadRequest, "Role binding name is required")
    }

    err := clientset.RbacV1().RoleBindings(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
    if err != nil {
        return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to delete role binding", err, "Failed to delete role binding in Kubernetes")
    }

    utils.Logger.Info("Role binding deleted successfully", zap.String("roleBindingName", name), zap.String("namespace", namespace))
    utils.LogAuditEvent(c.Request(), "delete", name, namespace)
    return c.JSON(http.StatusOK, map[string]string{"message": "Role binding deleted successfully"})
}

// RoleBindingDetailsHandler handles fetching detailed information about a specific role binding.
func RoleBindingDetailsHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		if !auth.HasPermission(username, "view_rolebinding_details") {
			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view role binding details")
		}

		roleBindingName := c.QueryParam("name")
		namespace := c.QueryParam("namespace")
		if namespace == "" {
			namespace = "default"
		}

		roleBinding, err := clientset.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBindingName, metav1.GetOptions{})
		if err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error fetching role binding details", err, "Failed to fetch role binding details")
		}

		utils.Logger.Info("Fetched role binding details", zap.String("roleBindingName", roleBindingName), zap.String("namespace", namespace))
		return c.JSON(http.StatusOK, roleBinding)
	}
}

// validateRoleBinding ensures that the role binding is valid.
func validateRoleBinding(roleBinding *rbacv1.RoleBinding) error {
	if roleBinding.Name == "" {
		return errors.New("role binding name is required")
	}
	if roleBinding.RoleRef.Name == "" {
		return errors.New("role reference name is required")
	}
	if len(roleBinding.Subjects) == 0 {
		return errors.New("at least one subject is required")
	}
	return nil
}
