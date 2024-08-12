package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RolesHandler returns an HTTP handler for managing roles
func RolesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch r.Method {
		case http.MethodGet:
			listRoles(w, clientset, namespace)
		case http.MethodPost:
			createRole(w, r, clientset, namespace)
		case http.MethodDelete:
			deleteRole(w, clientset, namespace, r.URL.Query().Get("name"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// listRoles lists all roles in the specified namespace
func listRoles(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace string) {
	roles, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles.Items)
}

// createRole creates a new role in the specified namespace
func createRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	var role rbacv1.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdRole, err := clientset.RbacV1().Roles(namespace).Create(context.TODO(), &role, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdRole)
}

// deleteRole deletes the role with the specified name in the specified namespace
func deleteRole(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace, name string) {
	if err := clientset.RbacV1().Roles(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		http.Error(w, "Failed to delete role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
