package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// roleDetailsHandler handles fetching detailed information about a specific role
func RoleDetailsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	// Get role name and namespace from query parameters
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getRoleDetails(w, r, clientset)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// getRoleDetails fetches detailed information about a specific role
func getRoleDetails(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	// Get the role name and namespace from the query parameters
	roleName := r.URL.Query().Get("roleName")
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	// Fetch the Role details
	role, err := clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the associated RoleBindings
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
