package server

import (
	"net/http"
	"os"

	"rbac/pkg/handlers/rbac"

	"github.com/labstack/echo/v4"
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

// RegisterRoutes registers all the routes for the server.
func RegisterRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	api := e.Group("/api")

	// Namespace routes
	api.GET("/namespaces", rbac.NamespacesHandler(clientset))
	api.POST("/namespaces", rbac.NamespacesHandler(clientset))
	api.DELETE("/namespaces", rbac.NamespacesHandler(clientset))

	// Role routes
	api.GET("/roles", rbac.RolesHandler(clientset))
	api.POST("/roles", rbac.RolesHandler(clientset))
	api.PUT("/roles", rbac.RolesHandler(clientset))
	api.DELETE("/roles", rbac.RolesHandler(clientset))
	api.GET("/roles/details", rbac.RoleDetailsHandler(clientset))

	// Role binding routes
	api.GET("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.POST("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.PUT("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.DELETE("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.GET("/rolebinding/details", rbac.RoleBindingDetailsHandler(clientset))

	// Cluster role routes
	api.GET("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.POST("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.PUT("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.DELETE("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.GET("/clusterroles/details", rbac.ClusterRoleDetailsHandler(clientset))

	// Cluster role binding routes
	api.GET("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.POST("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.PUT("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.DELETE("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.GET("/clusterrolebinding/details", rbac.ClusterRoleBindingDetailsHandler(clientset))

	// Service account routes
	api.GET("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.POST("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.DELETE("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.GET("/serviceaccount-details", rbac.ServiceAccountDetailsHandler(clientset))

	// Resource routes
	api.GET("/resources", rbac.APIResourcesHandler(clientset))

	// User routes
	api.GET("/users", rbac.UsersHandler(clientset))
	api.GET("/userroles", rbac.UserRolesHandler(clientset))

	// Group routes
	api.GET("/groups", rbac.GroupsHandler(clientset))
	api.GET("/groupdetails", rbac.GroupDetailsHandler(clientset))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Root URL handler
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Welcome to the Kubeberus"})
	})
}
