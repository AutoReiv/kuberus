package rbac

import (
	"context"
	"encoding/json"
	"net/http"
	"rbac/pkg/utils"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ClusterRolesHandler handles requests related to cluster roles.
func ClusterRolesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListClusterRoles(w, clientset)
		case http.MethodPost:
			handleCreateClusterRole(w, r, clientset)
		case http.MethodPut:
			handleUpdateClusterRole(w, r, clientset)
		case http.MethodDelete:
			handleDeleteClusterRole(w, r, clientset, r.URL.Query().Get("name"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// IsClusterRoleActive checks if a cluster role is active by looking for any cluster role bindings that reference it.
func IsClusterRoleActive(clientset *kubernetes.Clientset, clusterRoleName string) (bool, error) {
	// Check ClusterRoleBindings
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, crb := range clusterRoleBindings.Items {
		if crb.RoleRef.Name == clusterRoleName {
			return true, nil
		}
	}
	return false, nil
}

// ClusterRoleWithStatus represents a cluster role with its active status.
type ClusterRoleWithStatus struct {
	rbacv1.ClusterRole
	Active bool `json:"active"`
}

// handleListClusterRoles lists all cluster roles.
func handleListClusterRoles(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var clusterRolesWithStatus []ClusterRoleWithStatus
	for _, clusterRole := range clusterRoles.Items {
		active, err := IsClusterRoleActive(clientset, clusterRole.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clusterRolesWithStatus = append(clusterRolesWithStatus, ClusterRoleWithStatus{ClusterRole: clusterRole, Active: active})
	}

	utils.WriteJSON(w, clusterRolesWithStatus)
}

// handleCreateClusterRole creates a new cluster role.
func handleCreateClusterRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	var clusterRole rbacv1.ClusterRole
	if err := json.NewDecoder(r.Body).Decode(&clusterRole); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdClusterRole, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(), &clusterRole, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create cluster role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "create", clusterRole.Name, "cluster-wide")
	utils.WriteJSON(w, createdClusterRole)
}

// handleUpdateClusterRole updates an existing cluster role.
func handleUpdateClusterRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	var clusterRole rbacv1.ClusterRole
	if err := json.NewDecoder(r.Body).Decode(&clusterRole); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedClusterRole, err := clientset.RbacV1().ClusterRoles().Update(context.TODO(), &clusterRole, metav1.UpdateOptions{})
	if err != nil {
		http.Error(w, "Failed to update cluster role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "update", clusterRole.Name, "cluster-wide")
	utils.WriteJSON(w, updatedClusterRole)
}

// handleDeleteClusterRole deletes a cluster role by name.
func handleDeleteClusterRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, name string) {
	if name == "" {
		http.Error(w, "Cluster role name is required", http.StatusBadRequest)
		return
	}

	err := clientset.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete cluster role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "delete", name, "cluster-wide")
	w.WriteHeader(http.StatusNoContent)
}

// ClusterRoleDetailsHandler handles fetching detailed information about a specific cluster role.
func ClusterRoleDetailsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleGetClusterRoleDetails(w, r, clientset)
	}
}

// ClusterRoleDetailsResponse represents the detailed information about a cluster role.
type ClusterRoleDetailsResponse struct {
	ClusterRole         *rbacv1.ClusterRole         `json:"clusterRole"`
	ClusterRoleBindings []rbacv1.ClusterRoleBinding `json:"clusterRoleBindings"`
	Active              bool                        `json:"active"`
}

// handleGetClusterRoleDetails fetches detailed information about a specific cluster role.
func handleGetClusterRoleDetails(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	clusterRoleName := r.URL.Query().Get("clusterRoleName")
	if clusterRoleName == "" {
		http.Error(w, "Cluster role name is required", http.StatusBadRequest)
		return
	}

	clusterRole, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), clusterRoleName, metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	associatedBindings := filterClusterRoleBindings(clusterRoleBindings.Items, clusterRoleName)

	active, err := IsClusterRoleActive(clientset, clusterRoleName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ClusterRoleDetailsResponse{
		ClusterRole:         clusterRole,
		ClusterRoleBindings: associatedBindings,
		Active:              active,
	}

	utils.WriteJSON(w, response)
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

