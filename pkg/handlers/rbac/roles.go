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

// RolesHandler handles role-related requests.
func RolesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch r.Method {
		case http.MethodGet:
			handleGetRoles(w, clientset, namespace)
		case http.MethodPost:
			handleCreateRole(w, r, clientset, namespace)
		case http.MethodPut:
			handleUpdateRole(w, r, clientset, namespace)
		case http.MethodDelete:
			handleDeleteRole(w, clientset, namespace, r.URL.Query().Get("name"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleGetRoles handles listing roles in a specific namespace or across all namespaces.
func handleGetRoles(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace string) {
	if namespace == "all" {
		listAllNamespacesRoles(w, clientset)
	} else {
		listNamespaceRoles(w, clientset, namespace)
	}
}

// listNamespaceRoles lists roles in a specific namespace.
func listNamespaceRoles(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace string) {
	roles, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, roles)
}

// listAllNamespacesRoles lists roles across all namespaces.
func listAllNamespacesRoles(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, roles)
}

// handleCreateRole handles creating a new role in a specific namespace.
func handleCreateRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
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

	utils.WriteJSON(w, createdRole)
}

// handleUpdateRole handles updating an existing role in a specific namespace.
func handleUpdateRole(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	var role rbacv1.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedRole, err := clientset.RbacV1().Roles(namespace).Update(context.TODO(), &role, metav1.UpdateOptions{})
	if err != nil {
		http.Error(w, "Failed to update role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, updatedRole)
}

// handleDeleteRole handles deleting a role in a specific namespace.
func handleDeleteRole(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace, name string) {
	if name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	err := clientset.RbacV1().Roles(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RoleDetailsHandler handles fetching detailed information about a specific role.
func RoleDetailsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getRoleDetails(w, r, clientset)
	}
}

// getRoleDetails fetches detailed information about a specific role.
func getRoleDetails(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	roleName := r.URL.Query().Get("roleName")
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	role, err := clientset.RbacV1().Roles(namespace).Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	associatedBindings := filterRoleBindings(roleBindings.Items, roleName)

	response := RoleDetailsResponse{
		Role:         role,
		RoleBindings: associatedBindings,
	}

	utils.WriteJSON(w, response)
}

// filterRoleBindings filters role bindings associated with a specific role.
func filterRoleBindings(roleBindings []rbacv1.RoleBinding, roleName string) []rbacv1.RoleBinding {
	var associatedBindings []rbacv1.RoleBinding
	for _, rb := range roleBindings {
		if rb.RoleRef.Name == roleName {
			associatedBindings = append(associatedBindings, rb)
		}
	}
	return associatedBindings
}

// RoleDetailsResponse represents the detailed information about a role.
type RoleDetailsResponse struct {
	Role         *rbacv1.Role         `json:"role"`
	RoleBindings []rbacv1.RoleBinding `json:"roleBindings"`
}