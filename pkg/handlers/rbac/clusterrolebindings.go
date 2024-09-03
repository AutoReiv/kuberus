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

// ClusterRoleBindingsHandler handles requests related to cluster role bindings.
func ClusterRoleBindingsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListClusterRoleBindings(w, clientset)
		case http.MethodPost:
			handleCreateClusterRoleBinding(w, r, clientset)
		case http.MethodDelete:
			handleDeleteClusterRoleBinding(w, clientset, r.URL.Query().Get("name"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleListClusterRoleBindings lists all cluster role bindings.
func handleListClusterRoleBindings(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, clusterRoleBindings.Items)
}

// handleCreateClusterRoleBinding creates a new cluster role binding.
func handleCreateClusterRoleBinding(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	var clusterRoleBinding rbacv1.ClusterRoleBinding
	if err := json.NewDecoder(r.Body).Decode(&clusterRoleBinding); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdClusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), &clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create cluster role binding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, createdClusterRoleBinding)
}

// handleDeleteClusterRoleBinding deletes a cluster role binding by name.
func handleDeleteClusterRoleBinding(w http.ResponseWriter, clientset *kubernetes.Clientset, name string) {
	if name == "" {
		http.Error(w, "Cluster role binding name is required", http.StatusBadRequest)
		return
	}

	err := clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete cluster role binding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}