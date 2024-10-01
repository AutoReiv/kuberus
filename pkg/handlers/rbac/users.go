package rbac

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UserRolesHandler handles requests to show the roles or cluster roles a user has access to.
func UserRolesHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		userName := c.QueryParam("userName")
		if userName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "User name is required")
		}

		roleBindings, err := clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error listing role bindings: "+err.Error())
		}

		clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error listing cluster role bindings: "+err.Error())
		}

		userRoles := extractUserRoles(userName, roleBindings.Items, clusterRoleBindings.Items)
		return c.JSON(http.StatusOK, userRoles)
	}
}

// extractUserRoles extracts the roles and cluster roles a user has access to.
func extractUserRoles(username string, roleBindings []rbacv1.RoleBinding, clusterRoleBindings []rbacv1.ClusterRoleBinding) []string {
	var roles []string

	for _, rb := range roleBindings {
		for _, subject := range rb.Subjects {
			if subject.Kind == rbacv1.UserKind && subject.Name == username {
				roles = append(roles, rb.RoleRef.Name)
			}
		}
	}

	for _, crb := range clusterRoleBindings {
		for _, subject := range crb.Subjects {
			if subject.Kind == rbacv1.UserKind && subject.Name == username {
				roles = append(roles, crb.RoleRef.Name)
			}
		}
	}

	return roles
}

// UsersHandler handles requests to list all users from role bindings and cluster role bindings.
func UsersHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		roleBindings, err := clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error listing role bindings: "+err.Error())
		}

		clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error listing cluster role bindings: "+err.Error())
		}

		users := extractUsers(roleBindings.Items, clusterRoleBindings.Items)
		return c.JSON(http.StatusOK, users)
	}
}

// extractUsers extracts unique users from RoleBindings and ClusterRoleBindings.
func extractUsers(roleBindings []rbacv1.RoleBinding, clusterRoleBindings []rbacv1.ClusterRoleBinding) []string {
	userSet := make(map[string]struct{})

	for _, rb := range roleBindings {
		for _, subject := range rb.Subjects {
			if subject.Kind == rbacv1.UserKind {
				userSet[subject.Name] = struct{}{}
			}
		}
	}

	for _, crb := range clusterRoleBindings {
		for _, subject := range crb.Subjects {
			if subject.Kind == rbacv1.UserKind {
				userSet[subject.Name] = struct{}{}
			}
		}
	}

	users := make([]string, 0, len(userSet))
	for user := range userSet {
		users = append(users, user)
	}

	return users
}