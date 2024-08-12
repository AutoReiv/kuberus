package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// clusterRolesHandler handles listing cluster roles
func ClusterRolesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			listClusterRoles(w, clientset)
		case http.MethodPost:
			createClusterRole(w, r, clientset)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// listClusterRoles lists all cluster roles
func listClusterRoles(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusterRoles.Items)
}

// createClusterRole creates a new cluster role
func createClusterRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdClusterRole)
}
