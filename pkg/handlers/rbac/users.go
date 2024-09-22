package rbac

import (
	"context"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UserSource represents a user and their source.
type UserSource struct {
	Username string `json:"username"`
	Source   string `json:"source"`
}

// UsersHandler handles requests related to listing users.
func UsersHandler(clientset *kubernetes.Clientset) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Fetch users created by admin and OIDC users
		adminUsers, err := auth.GetAllUsers()
		if err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error retrieving users from database", err, "Failed to retrieve users from database")
		}

		// Fetch users from RoleBindings and ClusterRoleBindings
		roleBindings, err := clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing role bindings", err, "Failed to list role bindings")
		}

		clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error listing cluster role bindings", err, "Failed to list cluster role bindings")
		}

		k8sUsers := extractUsersFromBindings(roleBindings.Items, clusterRoleBindings.Items)

		// Combine the lists, ensuring no duplicates, and indicate the source
		userSet := make(map[string]string)
		for _, user := range adminUsers {
			userSet[user.Username] = user.Source
		}
		for _, user := range k8sUsers {
			if _, exists := userSet[user]; !exists {
				userSet[user] = "roleBinding"
			}
		}

		combinedUsers := make([]UserSource, 0, len(userSet))
		for user, source := range userSet {
			combinedUsers = append(combinedUsers, UserSource{Username: user, Source: source})
		}

		return c.JSON(http.StatusOK, combinedUsers)
	}
}

// extractUsersFromBindings extracts users from RoleBindings and ClusterRoleBindings.
func extractUsersFromBindings(roleBindings []rbacv1.RoleBinding, clusterRoleBindings []rbacv1.ClusterRoleBinding) []string {
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
