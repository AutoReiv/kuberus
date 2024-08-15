package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
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
		case http.MethodPut:
			editRole(w, r, clientset, namespace)
		case http.MethodDelete:
			deleteRole(w, clientset, namespace, r.URL.Query().Get("name"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// listRoles lists all roles in the specified namespace
func listRoles(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace string) {
	roles, err := clientset.RbacV1().Roles(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, roles.Items)
}

// createRole creates a new role in the specified namespace
func createRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	// Read the YAML payload from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Validate the YAML payload
	if err := validateKubernetesYAML(body); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Decode the validated YAML into a Role object
	var role rbacv1.Role
	if err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(body), 100).Decode(&role); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Create the role
	createdRole, err := clientset.RbacV1().Roles(namespace).Create(context.Background(), &role, metav1.CreateOptions{})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, createdRole)
}

// editRole edits the role with the specified name in the specified namespace
func editRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	// Read the YAML payload from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Validate the YAML payload
	if err := validateKubernetesYAML(body); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Decode the validated YAML into a Role object
	var role rbacv1.Role
	if err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(body), 100).Decode(&role); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	// Edit the role
	editedRole, err := clientset.RbacV1().Roles(namespace).Update(context.Background(), &role, metav1.UpdateOptions{})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, editedRole)
}

// deleteRole deletes the role with the specified name in the specified namespace
func deleteRole(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace, name string) {
	err := clientset.RbacV1().Roles(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleError sends an HTTP error response with the given error message and status code
func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

// respondWithJSON sends an HTTP response with the given data encoded as JSON
func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
