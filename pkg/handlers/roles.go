package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
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
			if namespace == "all" {
				listRolesAllNamespaces(w, clientset)
			} else {
				listRoles(w, clientset, namespace)
			}
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
	roles, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(roles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func listRolesAllNamespaces(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(roles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// createRole creates a new role in the specified namespace
func createRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate JSON content
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		http.Error(w, "Invalid JSON content: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convert JSON to Kubernetes object
	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer()
	obj, _, err := decoder.Decode(body, nil, nil)
	if err != nil {
		http.Error(w, "Invalid Kubernetes object: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the Kubernetes object using the client
	config, err := rest.InClusterConfig()
	if err != nil {
		http.Error(w, "Failed to get Kubernetes config: "+err.Error(), http.StatusInternalServerError)
		return
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		http.Error(w, "Failed to create Kubernetes client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Assuming the object is a Role, validate it
	role, ok := obj.(*rbacv1.Role)
	if !ok {
		http.Error(w, "Invalid Role object", http.StatusBadRequest)
		return
	}

	// Validate the role (this is a placeholder, actual validation logic may vary)
	if role.Name == "" || role.Namespace == "" {
		http.Error(w, "Role must have a name and namespace", http.StatusBadRequest)
		return
	}

	// Create the role in the specified namespace
	createdRole, err := clientset.RbacV1().Roles(namespace).Create(context.TODO(), &role, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created role
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdRole)
}

// editRole edits the role with the specified name in the specified namespace
func editRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the YAML content
	if err := validateKubernetesYAML(body); err != nil {
		http.Error(w, "Invalid YAML content: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal the JSON body into a Role object
	var role rbacv1.Role
	if err := json.Unmarshal(body, &role); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update the role in the specified namespace
	updatedRole, err := clientset.RbacV1().Roles(namespace).Update(context.TODO(), &role, metav1.UpdateOptions{})
	if err != nil {
		http.Error(w, "Failed to update role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated role
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRole)
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
