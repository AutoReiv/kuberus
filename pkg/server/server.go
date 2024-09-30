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
	e.GET("/auth/oidc/login", handlers.OIDCAuthHandler)
	e.GET("/auth/oidc/callback", handlers.OIDCCallbackHandler)
	// Admin OIDC configuration route
	e.POST("/admin/oidc/config", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(handlers.SetOIDCConfigHandler))
	// Admin Certificate Upload route
	uploadCertsHandler := handlers.NewUploadCertsHandler(clientset)
	e.POST("/admin/upload-certs", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(uploadCertsHandler.ServeHTTP))

	// User management routes
	e.POST("/admin/users", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(func(c echo.Context) error {
		return handlers.HandleCreateUser(c, clientset)
	}))
	e.DELETE("/admin/users", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(handlers.UserManagementHandler(clientset)))
	e.GET("/admin/users", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(handlers.UserManagementHandler(clientset)))
	e.PUT("/admin/users", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(handlers.UserManagementHandler(clientset)))
	// Protected API routes
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	// Namespace routes
	api.GET("/namespaces", middleware.AuthAndRBACMiddleware("list_namespaces", config.IsDevMode)(rbac.NamespacesHandler(clientset)))
	api.POST("/namespaces", middleware.AuthAndRBACMiddleware("create_namespace", config.IsDevMode)(rbac.NamespacesHandler(clientset)))
	api.DELETE("/namespaces", middleware.AuthAndRBACMiddleware("delete_namespace", config.IsDevMode)(rbac.NamespacesHandler(clientset)))

	// Role routes
	api.GET("/roles", middleware.AuthAndRBACMiddleware("list_roles", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.POST("/roles", middleware.AuthAndRBACMiddleware("create_role", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.PUT("/roles", middleware.AuthAndRBACMiddleware("update_role", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.DELETE("/roles", middleware.AuthAndRBACMiddleware("delete_role", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.GET("/roles/details", middleware.AuthAndRBACMiddleware("view_role_details", config.IsDevMode)(rbac.RoleDetailsHandler(clientset)))

	// Role bindings routes
	api.GET("/rolebindings", middleware.AuthAndRBACMiddleware("list_rolebindings", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.POST("/rolebindings", middleware.AuthAndRBACMiddleware("create_rolebinding", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.PUT("/rolebindings", middleware.AuthAndRBACMiddleware("update_rolebinding", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.DELETE("/rolebindings", middleware.AuthAndRBACMiddleware("delete_rolebinding", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.GET("/rolebinding/details", middleware.AuthAndRBACMiddleware("view_rolebinding_details", config.IsDevMode)(rbac.RoleBindingDetailsHandler(clientset)))

	// Cluster role routes
	api.GET("/clusterroles", middleware.AuthAndRBACMiddleware("list_clusterroles", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.POST("/clusterroles", middleware.AuthAndRBACMiddleware("create_clusterrole", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.PUT("/clusterroles", middleware.AuthAndRBACMiddleware("update_clusterrole", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.DELETE("/clusterroles", middleware.AuthAndRBACMiddleware("delete_clusterrole", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.GET("/clusterroles/details", middleware.AuthAndRBACMiddleware("view_clusterrole_details", config.IsDevMode)(rbac.ClusterRoleDetailsHandler(clientset)))

	// Cluster role bindings routes
	api.GET("/clusterrolebindings", middleware.AuthAndRBACMiddleware("list_clusterrolebindings", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.POST("/clusterrolebindings", middleware.AuthAndRBACMiddleware("create_clusterrolebinding", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.PUT("/clusterrolebindings", middleware.AuthAndRBACMiddleware("update_clusterrolebinding", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.DELETE("/clusterrolebindings", middleware.AuthAndRBACMiddleware("delete_clusterrolebinding", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.GET("/clusterrolebinding/details", middleware.AuthAndRBACMiddleware("view_clusterrolebinding_details", config.IsDevMode)(rbac.ClusterRoleBindingDetailsHandler(clientset)))

	// Resource routes
	api.GET("/resources", middleware.AuthAndRBACMiddleware("list_resources", config.IsDevMode)(rbac.APIResourcesHandler(clientset)))

	// Account routes
	api.GET("/serviceaccounts", middleware.AuthAndRBACMiddleware("list_serviceaccounts", config.IsDevMode)(rbac.ServiceAccountsHandler(clientset)))
	api.POST("/serviceaccounts", middleware.AuthAndRBACMiddleware("create_serviceaccount", config.IsDevMode)(rbac.ServiceAccountsHandler(clientset)))
	api.DELETE("/serviceaccounts", middleware.AuthAndRBACMiddleware("delete_serviceaccount", config.IsDevMode)(rbac.ServiceAccountsHandler(clientset)))
	api.GET("/serviceaccount-details", middleware.AuthAndRBACMiddleware("view_serviceaccount_details", config.IsDevMode)(rbac.ServiceAccountDetailsHandler(clientset)))

	// User routes
	api.GET("/users", middleware.AuthAndRBACMiddleware("list_users", config.IsDevMode)(rbac.UsersHandler(clientset)))
	api.GET("/user-details", middleware.AuthAndRBACMiddleware("view_user_details", config.IsDevMode)(rbac.UserDetailsHandler(clientset)))

	// Group routes
	api.GET("/groups", middleware.AuthAndRBACMiddleware("list_groups", config.IsDevMode)(rbac.GroupsHandler(clientset)))
	api.GET("/group-details", middleware.AuthAndRBACMiddleware("view_group_details", config.IsDevMode)(rbac.GroupDetailsHandler(clientset)))

	// User roles route
	api.GET("/user-roles", middleware.AuthAndRBACMiddleware("view_user_roles", config.IsDevMode)(rbac.UserRolesHandler(clientset)))

	// Audit logs route
	api.GET("/audit-logs", middleware.AuthAndRBACMiddleware("view_audit_logs", config.IsDevMode)(handlers.GetAuditLogsHandler))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Root URL handler
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Welcome to the RBAC Manager"})
	})

	// Simulation route
	api.POST("/simulate", middleware.AuthAndRBACMiddleware("simulate", config.IsDevMode)(rbac.SimulateHandler(clientset)))
}
