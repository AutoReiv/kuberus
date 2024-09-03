package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"rbac/pkg/handlers"
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"k8s.io/client-go/kubernetes"
)

// Config holds the configuration for the server.
type Config struct {
	Port      string
	IsDevMode bool
}

// NewConfig creates a new configuration with environment variables.
func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	isDevMode := os.Getenv("DEV_MODE") == "true"

	return &Config{Port: port, IsDevMode: isDevMode}
}

// NewServer creates a new HTTP server with the necessary routes and middleware.
func NewServer(clientset *kubernetes.Clientset, config *Config) *http.Server {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	registerRoutes(mux, clientset, config)

	// Apply middlewares
	handler := middleware.ApplyMiddlewares(mux, config.IsDevMode)

	// Create the HTTP server
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	handleGracefulShutdown(srv)

	return srv
}

// registerRoutes registers all the routes for the server.
func registerRoutes(mux *http.ServeMux, clientset *kubernetes.Clientset, config *Config) {
	// Admin account creation route
	mux.HandleFunc("/admin/create", handlers.CreateAdminHandler)

	// Authentication routes
	mux.HandleFunc("/auth/login", handlers.LoginHandler)
	// OIDC routes
	mux.HandleFunc("/auth/oidc/login", handlers.OIDCAuthHandler)
	mux.HandleFunc("/auth/oidc/callback", handlers.OIDCCallbackHandler)

	// Admin OIDC configuration route
	mux.Handle("/admin/oidc/config", middleware.AuthMiddleware(http.HandlerFunc(handlers.SetOIDCConfigHandler), config.IsDevMode))

	// Protected API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/namespaces", rbac.NamespacesHandler(clientset))
	apiMux.HandleFunc("/api/roles", rbac.RolesHandler(clientset))
	apiMux.HandleFunc("/api/roles/details", rbac.RoleDetailsHandler(clientset))
	apiMux.HandleFunc("/api/rolebindings", rbac.RoleBindingsHandler(clientset))
	apiMux.HandleFunc("/api/clusterroles", rbac.ClusterRolesHandler(clientset))
	apiMux.HandleFunc("/api/clusterroles/details", rbac.ClusterRoleDetailsHandler(clientset))
	apiMux.HandleFunc("/api/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	apiMux.HandleFunc("/api/resources", rbac.APIResourcesHandler(clientset))
	apiMux.HandleFunc("/api/serviceaccounts", rbac.ServiceAccountsHandler(clientset))

	mux.Handle("/api/", middleware.AuthMiddleware(apiMux, config.IsDevMode))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Root URL handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Welcome to the RBAC Manager"}`))
	})
}

// handleGracefulShutdown handles the graceful shutdown of the server.
func handleGracefulShutdown(srv *http.Server) {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()
}
