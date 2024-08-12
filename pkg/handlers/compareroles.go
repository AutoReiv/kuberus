package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CompareRolesHandler handles comparing roles
func CompareRolesHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			compareRoles(w, r, clientset)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// compareRoles compares two roles
func compareRoles(w http.ResponseWriter, r *http.Request, clientset *kubernetes.Clientset) {
	var roleNames struct {
		Role1 string `json:"role1"`
		Role2 string `json:"role2"`
	}
	if err := json.NewDecoder(r.Body).Decode(&roleNames); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	role1, err := clientset.RbacV1().Roles(r.URL.Query().Get("namespace")).Get(context.TODO(), roleNames.Role1, metav1.GetOptions{})
	if err != nil {
		http.Error(w, "Failed to get role1: "+err.Error(), http.StatusInternalServerError)
		return
	}

	role2, err := clientset.RbacV1().Roles(r.URL.Query().Get("namespace")).Get(context.TODO(), roleNames.Role2, metav1.GetOptions{})
	if err != nil {
		http.Error(w, "Failed to get role2: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Compare the two roles
	// ...

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Role1 interface{} `json:"role1"`
		Role2 interface{} `json:"role2"`
	}{role1, role2})
}
