package server

import (
	"net/http"
	"os"

	"rbac/pkg/handlers"
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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

	// CORS middleware
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Admin account creation route
	e.POST("/admin/create", handlers.CreateAdminHandler)

	// Authentication routes
	e.POST("/auth/login", handlers.LoginHandler)
	// OIDC routes
	e.GET("/auth/oidc/login", handlers.OIDCAuthHandler)
	e.GET("/auth/oidc/callback", handlers.OIDCCallbackHandler)

	// Admin OIDC configuration route
	e.POST("/admin/oidc/config", middleware.JWTMiddleware()(handlers.SetOIDCConfigHandler))
	// Admin Certificate Upload route
	uploadCertsHandler := handlers.NewUploadCertsHandler(clientset)
	e.POST("/admin/upload-certs", middleware.JWTMiddleware()(uploadCertsHandler.ServeHTTP))

	// User management routes
	e.POST("/admin/users", middleware.JWTMiddleware()(func(c echo.Context) error {
		return handlers.HandleCreateUser(c, clientset)
	}))
	e.DELETE("/admin/users", middleware.JWTMiddleware()(handlers.UserManagementHandler(clientset)))
	e.GET("/admin/users", middleware.JWTMiddleware()(handlers.UserManagementHandler(clientset)))
	e.PUT("/admin/users", middleware.JWTMiddleware()(handlers.UserManagementHandler(clientset)))

	// Protected API routes
	api := e.Group("/api", middleware.JWTMiddleware())
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

	// Role bindings routes
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

	// Cluster role bindings routes
	api.GET("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.POST("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.PUT("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.DELETE("/clusterrolebindings", rbac.ClusterRoleBindingsHandler(clientset))
	api.GET("/clusterrolebinding/details", rbac.ClusterRoleBindingDetailsHandler(clientset))

	// Resource routes
	api.GET("/resources", rbac.APIResourcesHandler(clientset))

	// Account routes
	api.GET("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.POST("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.DELETE("/serviceaccounts", rbac.ServiceAccountsHandler(clientset))
	api.GET("/serviceaccount-details", rbac.ServiceAccountDetailsHandler(clientset))

	// User routes
	api.GET("/users", rbac.UsersHandler(clientset))
	api.GET("/user-details", rbac.UserDetailsHandler(clientset))

	// Group routes
	api.GET("/groups", rbac.GroupsHandler(clientset))
	api.GET("/group-details", rbac.GroupDetailsHandler(clientset))

	// User roles route
	api.GET("/user-roles", middleware.JWTMiddleware()(rbac.UserRolesHandler(clientset)))

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

	// Simulation route
	api.POST("/simulate", rbac.SimulateHandler(clientset))
}
