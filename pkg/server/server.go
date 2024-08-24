package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rbac/pkg/auth"
	"rbac/pkg/handlers"
	"rbac/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

func NewServer(clientset *kubernetes.Clientset) *http.Server {
	r := gin.Default()

	// Admin account creation route
	r.POST("/admin/create", handlers.CreateAdminHandler)

	// Register the HTTP handlers
	r.POST("/login", handlers.LoginHandler)
	r.POST("/register", handlers.RegisterHandler)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware)
	api.GET("/namespaces", handlers.NamespacesHandler(clientset))
	api.GET("/roles", handlers.RolesHandler(clientset))
	api.GET("/roles/details", handlers.RoleDetailsHandler(clientset))
	api.GET("/clusterroles", handlers.ClusterRolesHandler(clientset))

	// Add health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Configure the OIDC provider
	auth.ConfigureOIDCProvider()

	// OAuth routes
	r.GET("/auth/login", handlers.OAuthLoginHandler)
	r.GET("/auth/callback", handlers.OAuthCallbackHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
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

	return srv
}
