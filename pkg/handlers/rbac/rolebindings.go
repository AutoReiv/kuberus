package rbac

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"rbac/pkg/utils"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RoleBindingsHandler handles role binding-related requests.
func RoleBindingsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch r.Method {
		case http.MethodGet:
			handleListRoleBindings(w, clientset, namespace)
		case http.MethodPost:
			handleCreateRoleBinding(w, r, clientset, namespace)
		case http.MethodPut:
			handleUpdateRoleBinding(w, r, clientset, namespace)
		case http.MethodDelete:
			handleDeleteRoleBinding(w, r, clientset, namespace, r.URL.Query().Get("name"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleListRoleBindings lists all role bindings in a specific namespace.
func handleListRoleBindings(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace string) {
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, roleBindings.Items)
}

// handleCreateRoleBinding creates a new role binding in a specific namespace.
func handleCreateRoleBinding(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	var roleBinding rbacv1.RoleBinding
	if err := json.NewDecoder(r.Body).Decode(&roleBinding); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateRoleBinding(&roleBinding); err != nil {
		http.Error(w, "Invalid role binding: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdRoleBinding, err := clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), &roleBinding, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create role binding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "create", roleBinding.Name, namespace)
	utils.WriteJSON(w, createdRoleBinding)
}

// handleUpdateRoleBinding updates an existing role binding in a specific namespace.
func handleUpdateRoleBinding(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	var roleBinding rbacv1.RoleBinding
	if err := json.NewDecoder(r.Body).Decode(&roleBinding); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateRoleBinding(&roleBinding); err != nil {
		http.Error(w, "Invalid role binding: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedRoleBinding, err := clientset.RbacV1().RoleBindings(namespace).Update(context.TODO(), &roleBinding, metav1.UpdateOptions{})
	if err != nil {
		http.Error(w, "Failed to update role binding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "update", roleBinding.Name, namespace)
	utils.WriteJSON(w, updatedRoleBinding)
}

// handleDeleteRoleBinding deletes a role binding in a specific namespace.
func handleDeleteRoleBinding(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace, name string) {
	if name == "" {
		http.Error(w, "Role binding name is required", http.StatusBadRequest)
		return
	}

	err := clientset.RbacV1().RoleBindings(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete role binding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "delete", name, namespace)
	w.WriteHeader(http.StatusNoContent)
}

// RoleBindingDetailsHandler handles fetching detailed information about a specific role binding.
func RoleBindingDetailsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleBindingName := r.URL.Query().Get("name")
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default"
		}

		roleBinding, err := clientset.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBindingName, metav1.GetOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, roleBinding)
	}
}

// validateRoleBinding ensures that the role binding is valid.
func validateRoleBinding(roleBinding *rbacv1.RoleBinding) error {
	if roleBinding.Name == "" {
		return errors.New("role binding name is required")
	}
	if roleBinding.RoleRef.Name == "" {
		return errors.New("role reference name is required")
	}
	if len(roleBinding.Subjects) == 0 {
		return errors.New("at least one subject is required")
	}
	return nil
}