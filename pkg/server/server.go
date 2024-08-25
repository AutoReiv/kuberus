package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"rbac/pkg/auth"
	"rbac/pkg/handlers"
	"rbac/pkg/middleware"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

// Config holds the configuration for the server.
type Config struct {
	Port string
}

// NewConfig creates a new configuration with environment variables.
func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{Port: port}
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
	registerRoutes(r, clientset)

	// Configure the OIDC provider
	auth.ConfigureOIDCProvider()

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
func registerRoutes(r *gin.Engine, clientset *kubernetes.Clientset) {
	// Admin account creation route
	r.POST("/admin/create", handlers.CreateAdminHandler)

	// Authentication routes
	r.POST("/login", handlers.LoginHandler)
	r.GET("/auth/login", handlers.OAuthLoginHandler)
	r.GET("/auth/callback", handlers.OAuthCallbackHandler)

	// Protected API routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	api.GET("/namespaces", handlers.NamespacesHandler(clientset))
	api.GET("/roles", handlers.RolesHandler(clientset))
	api.GET("/roles/details", handlers.RoleDetailsHandler(clientset))
	api.GET("/clusterroles", handlers.ClusterRolesHandler(clientset))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
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
