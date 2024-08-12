package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// clusterRoleBindingsHandler handles listing cluster role bindings
func ClusterRoleBindingsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			listClusterRoleBindings(w, clientset)
		case http.MethodPost:
			createClusterRoleBindings(w, r, clientset)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// listClusterRoleBindings lists all cluster role bindings
func listClusterRoleBindings(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusterRoleBindings.Items)
}

// createClusterRoleBindings creates a new cluster role binding
func createClusterRoleBindings(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdClusterRoleBinding)
}
