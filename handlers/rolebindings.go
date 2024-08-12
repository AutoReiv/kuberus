package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// roleBindingsHandler handles listing role bindings
func RoleBindingsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default" // Default to "default" namespace
		}

		switch r.Method {
		case http.MethodGet:
			listRoleBindings(w, clientset, namespace)
		case http.MethodPost:
			createRoleBindings(w, r, clientset, namespace)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// listRoleBindings lists all role bindings in the "default" namespace
func listRoleBindings(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace string) {
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roleBindings.Items)
}

// createRoleBindings creates a new role binding in the specified namespace
func createRoleBindings(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	var roleBinding rbacv1.RoleBinding
	if err := json.NewDecoder(r.Body).Decode(&roleBinding); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdRoleBinding, err := clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), &roleBinding, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create role binding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdRoleBinding)
}

