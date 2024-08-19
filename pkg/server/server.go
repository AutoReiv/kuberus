package server

import (
	"log"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/handlers"

	"k8s.io/client-go/kubernetes"
)

// NewServer initializes a new HTTP server
func NewServer(clientset *kubernetes.Clientset) *http.Server {
	provider, verifier, err := auth.InitOIDC()
	if err != nil {
		log.Fatalf("Failed to initialize OIDC: %v", err)
	}

	// Log provider information
	log.Printf("OIDC provider initialized: %v", provider)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the HTTP handlers
	mux.Handle("/roles", auth.AuthMiddleware(verifier, handlers.RolesHandler(clientset)))
	mux.HandleFunc("/api/namespaces", handlers.NamespacesHandler(clientset))
	mux.HandleFunc("/api/roles", handlers.RolesHandler(clientset))
	mux.HandleFunc("/api/roles/details", handlers.RoleDetailsHandler(clientset))
	mux.HandleFunc("/api/roles/compare", handlers.CompareRolesHandler(clientset))
	mux.HandleFunc("/api/clusterroles", handlers.ClusterRolesHandler(clientset))

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
