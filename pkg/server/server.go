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

	handleGracefulShutdown(srv)
	return srv
}

func StartServer(srv *http.Server, config *Config) error {
	certFile := "certs/tls.crt"
	keyFile := "certs/tls.key"

	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			return srv.ListenAndServeTLS(certFile, keyFile)
		}
	}

	return srv.ListenAndServe()
}

// registerRoutes registers all the routes for the server.
func registerRoutes(mux *http.ServeMux, clientset *kubernetes.Clientset, config *Config) {
	// Admin account creation route
	mux.Handle("/admin/create", http.HandlerFunc(handlers.CreateAdminHandler))

	// Authentication routes
	mux.Handle("/auth/login", http.HandlerFunc(handlers.LoginHandler))
	// OIDC routes
	mux.Handle("/auth/oidc/login", http.HandlerFunc(handlers.OIDCAuthHandler))
	mux.Handle("/auth/oidc/callback", http.HandlerFunc(handlers.OIDCCallbackHandler))

	// Admin OIDC configuration route
	mux.Handle("/admin/oidc/config", middleware.AuthMiddleware(http.HandlerFunc(handlers.SetOIDCConfigHandler), config.IsDevMode))
	// Admin Certificate Upload route
	uploadCertsHandler := handlers.NewUploadCertsHandler(clientset)
	mux.Handle("/admin/upload-certs", middleware.AuthMiddleware(uploadCertsHandler, config.IsDevMode))

	// Protected API routes
	apiMux := http.NewServeMux()
	apiMux.Handle("/api/namespaces", http.HandlerFunc(rbac.NamespacesHandler(clientset)))
	apiMux.Handle("/api/roles", http.HandlerFunc(rbac.RolesHandler(clientset)))
	apiMux.Handle("/api/roles/details", http.HandlerFunc(rbac.RoleDetailsHandler(clientset)))
	apiMux.Handle("/api/rolebindings", http.HandlerFunc(rbac.RoleBindingsHandler(clientset)))
	apiMux.Handle("/api/clusterroles", http.HandlerFunc(rbac.ClusterRolesHandler(clientset)))
	apiMux.Handle("/api/clusterroles/details", http.HandlerFunc(rbac.ClusterRoleDetailsHandler(clientset)))
	apiMux.Handle("/api/clusterrolebindings", http.HandlerFunc(rbac.ClusterRoleBindingsHandler(clientset)))
	apiMux.Handle("/api/resources", http.HandlerFunc(rbac.APIResourcesHandler(clientset)))
	apiMux.Handle("/api/serviceaccounts", http.HandlerFunc(rbac.ServiceAccountsHandler(clientset)))
	apiMux.Handle("/api/serviceaccount-details", http.HandlerFunc(rbac.ServiceAccountDetailsHandler(clientset)))
	apiMux.Handle("/api/users", http.HandlerFunc(rbac.UsersHandler(clientset)))
	apiMux.Handle("/api/user-details", http.HandlerFunc(rbac.UserDetailsHandler(clientset)))
	apiMux.Handle("/api/groups", http.HandlerFunc(rbac.GroupsHandler(clientset)))
	apiMux.Handle("/api/group-details", http.HandlerFunc(rbac.GroupDetailsHandler(clientset)))

	mux.Handle("/api/", middleware.AuthMiddleware(apiMux, config.IsDevMode))

	// Health check endpoint
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Root URL handler
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Welcome to the RBAC Manager"}`))
	}))
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
