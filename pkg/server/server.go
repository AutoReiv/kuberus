// pkg/server/server.go
package server

import (
	"net/http"
	"os"

	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

// Config holds the configuration for the server.
type Config struct {
	Port      string
	IsDevMode bool
	JWTSecret []byte
}

// NewConfig creates a new configuration with environment variables.
func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	isDevMode := os.Getenv("DEV_MODE") == "true"
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	return &Config{Port: port, IsDevMode: isDevMode, JWTSecret: jwtSecret}
}

// RegisterRoutes registers all the routes for the server.
func RegisterRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	// Apply middlewares
	middleware.ApplyMiddlewares(e, config.IsDevMode)

	// Register routes
	registerAdminRoutes(e, clientset, config)
	registerAuthRoutes(e)
	registerNamespaceRoutes(e, clientset, config)
	registerRoleRoutes(e, clientset, config)
	registerRoleBindingRoutes(e, clientset, config)
	registerClusterRoleRoutes(e, clientset, config)
	registerClusterRoleBindingRoutes(e, clientset, config)
	registerResourceRoutes(e, clientset, config)
	registerServiceAccountRoutes(e, clientset, config)
	registerAuditRoutes(e, config)
	registerAdminUserRoutes(e, clientset) // Register user routes under /admin/users

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Root URL handler
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Welcome to the RBAC Manager"})
	})
}