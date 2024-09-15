package server

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
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

// RegisterRoutes registers all the routes for the server.
func RegisterRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	// Admin account creation route
	e.POST("/admin/create", handlers.CreateAdminHandler)

	// Authentication routes
	e.POST("/auth/login", handlers.LoginHandler)
	// OIDC routes
	e.GET("/auth/oidc/login", handlers.OIDCAuthHandler)
	e.GET("/auth/oidc/callback", handlers.OIDCCallbackHandler)

	// Admin OIDC configuration route
	e.POST("/admin/oidc/config", middleware.AuthMiddleware(handlers.SetOIDCConfigHandler))
	// Admin Certificate Upload route
	uploadCertsHandler := handlers.NewUploadCertsHandler(clientset)
	e.POST("/admin/upload-certs", middleware.AuthMiddleware(uploadCertsHandler.ServeHTTP))

	// User management routes
	e.POST("/admin/users", middleware.AuthMiddleware(handlers.UserManagementHandler(clientset)))

	// Protected API routes
	api := e.Group("/api")
	// Namespace routes
	api.GET("/namespaces", rbac.NamespacesHandler(clientset))
	// Role routes
	api.GET("/roles", rbac.RolesHandler(clientset))
	api.POST("/roles", rbac.RolesHandler(clientset))
	api.PUT("/roles", rbac.RolesHandler(clientset))
	api.DELETE("/roles", rbac.RolesHandler(clientset))
	api.GET("/roles/details", rbac.RoleDetailsHandler(clientset))
	// Role bindings routes
	api.GET("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.POST("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.PUT("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.DELETE("/rolebindings", rbac.RoleBindingsHandler(clientset))
	api.GET("/rolebinding/details", rbac.RoleBindingDetailsHandler(clientset))
	// Cluster role routes
	api.GET("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.GET("/clusterroles/details", rbac.ClusterRoleDetailsHandler(clientset))
	// Cluster role bindings routes
	api.GET("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.GET("/clusterrolebinding/details", rbac.ClusterRoleBindingDetailsHandler(clientset))
	// Resource routes
	api.GET("/resources", rbac.APIResourcesHandler(clientset))
	// Account routes
	api.GET("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.GET("/serviceaccount-details", rbac.ServiceAccountDetailsHandler(clientset))
	api.GET("/users", rbac.UsersHandler(clientset))
	api.GET("/user-details", rbac.UserDetailsHandler(clientset))
	api.GET("/groups", rbac.GroupsHandler(clientset))
	api.GET("/group-details", rbac.GroupDetailsHandler(clientset))
	// Audit logs route
	api.GET("/audit-logs", handlers.GetAuditLogsHandler)	// Cluster role routes
	api.GET("/clusterroles", rbac.ClusterRolesHandler(clientset))
	api.GET("/clusterroles/details", rbac.ClusterRoleDetailsHandler(clientset))
	// Cluster role bindings routes
	api.GET("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.GET("/clusterrolebinding/details", rbac.ClusterRoleBindingDetailsHandler(clientset))
	// Resource routes
	api.GET("/resources", rbac.APIResourcesHandler(clientset))
	// Account routes
	api.GET("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.GET("/serviceaccount-details", rbac.ServiceAccountDetailsHandler(clientset))
	api.GET("/users", rbac.UsersHandler(clientset))
	api.GET("/user-details", rbac.UserDetailsHandler(clientset))
	api.GET("/groups", rbac.GroupsHandler(clientset))
	api.GET("/group-details", rbac.GroupDetailsHandler(clientset))
	// Audit logs route
	api.GET("/audit-logs", handlers.GetAuditLogsHandler)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Root URL handler
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Welcome to the RBAC Manager"})
	})
}