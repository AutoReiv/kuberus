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

// ServiceAccountsHandler handles requests related to service accounts.
func ServiceAccountsHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			namespace = "default"
		}

		switch r.Method {
		case http.MethodGet:
			handleListServiceAccounts(w, clientset, namespace)
		case http.MethodPost:
			handleCreateServiceAccount(w, r, clientset, namespace)
		case http.MethodDelete:
			handleDeleteServiceAccount(w, clientset, namespace, r.URL.Query().Get("name"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// handleListServiceAccounts lists all service accounts in a specific namespace.
func handleListServiceAccounts(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace string) {
	serviceAccounts, err := clientset.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, serviceAccounts.Items)
}

// handleCreateServiceAccount creates a new service account in a specific namespace.
func handleCreateServiceAccount(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset, namespace string) {
	var serviceAccount corev1.ServiceAccount
	if err := json.NewDecoder(r.Body).Decode(&serviceAccount); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdServiceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &serviceAccount, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "Failed to create service account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, createdServiceAccount)
}

// handleDeleteServiceAccount deletes a service account in a specific namespace.
func handleDeleteServiceAccount(w http.ResponseWriter, clientset *kubernetes.Clientset, namespace, name string) {
	if name == "" {
		http.Error(w, "Service account name is required", http.StatusBadRequest)
		return
	}

	err := clientset.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "Failed to delete service account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}