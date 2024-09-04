package rbac

import (
	"context"
	"net/http"
	"rbac/pkg/utils"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// UserDetailsResponse represents the detailed information about a user.
type UserDetailsResponse struct {
	UserName            string                      `json:"userName"`
	RoleBindings        []rbacv1.RoleBinding        `json:"roleBindings"`
	ClusterRoleBindings []rbacv1.ClusterRoleBinding `json:"clusterRoleBindings"`
	ClusterRoles        []rbacv1.ClusterRole        `json:"clusterRoles"`
}

// UserDetailsHandler handles requests for detailed information about a specific user.
func UserDetailsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userName := r.URL.Query().Get("userName")
		if userName == "" {
			http.Error(w, "User name is required", http.StatusBadRequest)
			return
		}

		roleBindings, err := clientset.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userDetails := extractUserDetails(userName, roleBindings.Items, clusterRoleBindings.Items, clusterRoles.Items)
		utils.WriteJSON(w, userDetails)
	}
}
		// extractUserDetails extracts detailed information about a specific user.
		func extractUserDetails(userName string, roleBindings []rbacv1.RoleBinding, clusterRoleBindings []rbacv1.ClusterRoleBinding, clusterRoles []rbacv1.ClusterRole) UserDetailsResponse {
			var userRoleBindings []rbacv1.RoleBinding
			var userClusterRoleBindings []rbacv1.ClusterRoleBinding
			var userClusterRoles []rbacv1.ClusterRole

			for _, rb := range roleBindings {
				for _, subject := range rb.Subjects {
					if subject.Kind == rbacv1.UserKind && subject.Name == userName {
						userRoleBindings = append(userRoleBindings, rb)
					}
				}
			}

			for _, crb := range clusterRoleBindings {
				for _, subject := range crb.Subjects {
					if subject.Kind == rbacv1.UserKind && subject.Name == userName {
						userClusterRoleBindings = append(userClusterRoleBindings, crb)
					}
				}
			}

			// Collect ClusterRoles associated with the user's ClusterRoleBindings
			for _, crb := range userClusterRoleBindings {
				for _, cr := range clusterRoles {
					if cr.Name == crb.RoleRef.Name {
						userClusterRoles = append(userClusterRoles, cr)
					}
				}
			}

			return UserDetailsResponse{
				UserName:            userName,
				RoleBindings:        userRoleBindings,
				ClusterRoleBindings: userClusterRoleBindings,
				ClusterRoles:        userClusterRoles,
			}
		}
