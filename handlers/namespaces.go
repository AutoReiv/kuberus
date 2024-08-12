package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespacesHandler handles listing namespaces
func NamespacesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listNamespaces(w, clientset)
		case http.MethodPost:
			createNamespace(w, r, clientset)
		case http.MethodDelete:
			deleteNamespace(w, clientset, r.URL.Query().Get("name"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// listNamespaces lists all namespaces
func listNamespaces(w http.ResponseWriter, clientset *kubernetes.Clientset) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(namespaces.Items)
}

// createNamespace creates a new namespace
func createNamespace(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdNamespace)
}

// deleteNamespace deletes a namespace
func deleteNamespace(w http.ResponseWriter, clientset *kubernetes.Clientset, name string) {
	if err := clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		http.Error(w, "Failed to delete namespace: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
