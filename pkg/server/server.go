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

	"github.com/gin-gonic/gin"
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
	// Create a new Gin router
	r := gin.New()

	// Use Gin's logger and recovery middleware for better logging and error handling
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Secure the server with secure headers
	r.Use(middleware.SecureHeaders())

	// Register routes
	registerRoutes(r, clientset, config)

	// Create the HTTP server
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	handleGracefulShutdown(srv)

	return srv
}

// registerRoutes registers all the routes for the server.
func registerRoutes(r *gin.Engine, clientset *kubernetes.Clientset, config *Config) {
	// Admin account creation route
	r.POST("/admin/create", handlers.CreateAdminHandler)

	// Authentication routes
	auth := r.Group("/auth")
	auth.POST("/login", handlers.LoginHandler)
	// OIDC routes
	auth.GET("/oidc/login", handlers.OIDCAuthHandler)
	auth.GET("/oidc/callback", handlers.OIDCCallbackHandler)

	// Admin OIDC configuration route
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(config.IsDevMode))
	admin.POST("/oidc/config", handlers.SetOIDCConfigHandler)

	// Protected API routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(config.IsDevMode))
	api.GET("/namespaces", rbac.NamespacesHandler(clientset))
	api.GET("/roles", rbac.RolesHandler(clientset))
	api.GET("/roles/details", rbac.RoleDetailsHandler(clientset))
	api.POST("/roles", rbac.RolesHandler(clientset))
	api.PUT("/roles", rbac.RolesHandler(clientset))
	api.GET("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.GET("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.GET("/clusterroles/details", rbac.ClusterRoleDetailsHandler(clientset))
	api.POST("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.GET("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.GET("/resources", rbac.APIResourcesHandler(clientset))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Root URL handler
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the RBAC Manager"})
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
