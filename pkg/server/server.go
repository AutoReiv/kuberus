package server

import (
	"net/http"
	"rbac/pkg/handlers"
	"rbac/pkg/middleware"

	"k8s.io/client-go/kubernetes"
)

// NewServer initializes a new HTTP server
func NewServer(clientset *kubernetes.Clientset) *http.Server {

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the HTTP handlers
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	// Protected routes
	mux.Handle("/api/namespaces", middleware.AuthMiddleware(http.HandlerFunc(handlers.NamespacesHandler(clientset))))
	mux.Handle("/api/roles", middleware.AuthMiddleware(http.HandlerFunc(handlers.RolesHandler(clientset))))
	mux.Handle("/api/roles/details", middleware.AuthMiddleware(http.HandlerFunc(handlers.RoleDetailsHandler(clientset))))
	mux.Handle("/api/roles/compare", middleware.AuthMiddleware(http.HandlerFunc(handlers.CompareRolesHandler(clientset))))
	mux.Handle("/api/clusterroles", middleware.AuthMiddleware(http.HandlerFunc(handlers.ClusterRolesHandler(clientset))))

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
