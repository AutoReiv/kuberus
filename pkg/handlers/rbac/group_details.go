package rbac

import (
	"context"
	"net/http"
	"rbac/pkg/utils"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GroupDetailsResponse represents the detailed information about a group.
type GroupDetailsResponse struct {
	GroupName           string                      `json:"groupName"`
	RoleBindings        []rbacv1.RoleBinding        `json:"roleBindings"`
	ClusterRoleBindings []rbacv1.ClusterRoleBinding `json:"clusterRoleBindings"`
	ClusterRoles        []rbacv1.ClusterRole        `json:"clusterRoles"`
}

// GroupDetailsHandler handles requests for detailed information about a specific group.
func GroupDetailsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupName := r.URL.Query().Get("groupName")
		if groupName == "" {
			http.Error(w, "Group name is required", http.StatusBadRequest)
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

		groupDetails := extractGroupDetails(groupName, roleBindings.Items, clusterRoleBindings.Items, clusterRoles.Items)
		utils.WriteJSON(w, groupDetails)
	}
}

// extractGroupDetails extracts detailed information about a specific group.
func extractGroupDetails(groupName string, roleBindings []rbacv1.RoleBinding, clusterRoleBindings []rbacv1.ClusterRoleBinding, clusterRoles []rbacv1.ClusterRole) GroupDetailsResponse {
	var groupRoleBindings []rbacv1.RoleBinding
	var groupClusterRoleBindings []rbacv1.ClusterRoleBinding
	var groupClusterRoles []rbacv1.ClusterRole

	for _, rb := range roleBindings {
		for _, subject := range rb.Subjects {
			if subject.Kind == rbacv1.GroupKind && subject.Name == groupName {
				groupRoleBindings = append(groupRoleBindings, rb)
			}
		}
	}

	for _, crb := range clusterRoleBindings {
		for _, subject := range crb.Subjects {
			if subject.Kind == rbacv1.GroupKind && subject.Name == groupName {
				groupClusterRoleBindings = append(groupClusterRoleBindings, crb)
			}
		}
	}

	// Collect ClusterRoles associated with the group's ClusterRoleBindings
	for _, crb := range groupClusterRoleBindings {
		for _, cr := range clusterRoles {
			if cr.Name == crb.RoleRef.Name {
				groupClusterRoles = append(groupClusterRoles, cr)
			}
		}
	}

	return GroupDetailsResponse{
		GroupName:           groupName,
		RoleBindings:        groupRoleBindings,
		ClusterRoleBindings: groupClusterRoleBindings,
		ClusterRoles:        groupClusterRoles,
	}
}