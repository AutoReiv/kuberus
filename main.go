package main

import (
	"log"
	"net/http"
	"rbac/handlers"
	"rbac/kubernetes"
)

func main() {
	// Initialize the Kubernetes client
	clientset, err := kubernetes.NewClientset()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Register the HTTP handlers
	http.HandleFunc("/api/namespaces", handlers.NamespacesHandler(clientset))
	http.HandleFunc("/api/roles", handlers.RolesHandler(clientset))
	http.HandleFunc("/api/roles/details", handlers.RoleDetailsHandler(clientset))
	http.HandleFunc("/api/roles/compare", handlers.CompareRolesHandler(clientset))
	http.HandleFunc("/api/clusterroles", handlers.ClusterRolesHandler(clientset))
	http.HandleFunc("/api/rolebindings", handlers.RoleBindingsHandler(clientset))
	http.HandleFunc("/api/clusterrolebindings", handlers.ClusterRoleBindingsHandler(clientset))

	// Start the server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
