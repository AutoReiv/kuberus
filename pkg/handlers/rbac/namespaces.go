package rbac

import (
	"context"
	"encoding/json"
	"net/http"
	"rbac/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespacesHandler handles requests related to namespaces.
func NamespacesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListNamespaces(w, clientset)
		case http.MethodPost:
			handleCreateNamespace(w, r, clientset)
		case http.MethodDelete:
			handleDeleteNamespace(w, clientset, r.URL.Query().Get("name"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// handleListNamespaces lists all namespaces.
func handleListNamespaces(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, namespaces.Items)
}

// handleCreateNamespace creates a new namespace.
func handleCreateNamespace(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	var namespace corev1.Namespace
	if err := json.NewDecoder(r.Body).Decode(&namespace); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdNamespace, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create namespace: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, createdNamespace)
}

// handleDeleteNamespace deletes a namespace by name.
func handleDeleteNamespace(w http.ResponseWriter, clientset *kubernetes.Clientset, name string) {
	if name == "" {
		http.Error(w, "Namespace name is required", http.StatusBadRequest)
		return
	}

	err := clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete namespace: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}