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
	api.GET("/namespaces", middleware.RBACMiddleware("list_namespaces")(rbac.NamespacesHandler(clientset)))
	api.POST("/namespaces", middleware.RBACMiddleware("create_namespace")(rbac.NamespacesHandler(clientset)))
	api.DELETE("/namespaces", middleware.RBACMiddleware("delete_namespace")(rbac.NamespacesHandler(clientset)))

	// Role routes
	api.GET("/roles", middleware.RBACMiddleware("list_roles")(rbac.RolesHandler(clientset)))
	api.POST("/roles", middleware.RBACMiddleware("create_role")(rbac.RolesHandler(clientset)))
	api.PUT("/roles", middleware.RBACMiddleware("update_role")(rbac.RolesHandler(clientset)))
	api.DELETE("/roles", middleware.RBACMiddleware("delete_role")(rbac.RolesHandler(clientset)))
	api.GET("/roles/details", middleware.RBACMiddleware("view_role_details")(rbac.RoleDetailsHandler(clientset)))

	// Role bindings routes
	api.GET("/rolebindings", middleware.RBACMiddleware("list_rolebindings")(rbac.RoleBindingsHandler(clientset)))
	api.POST("/rolebindings", middleware.RBACMiddleware("create_rolebinding")(rbac.RoleBindingsHandler(clientset)))
	api.PUT("/rolebindings", middleware.RBACMiddleware("update_rolebinding")(rbac.RoleBindingsHandler(clientset)))
	api.DELETE("/rolebindings", middleware.RBACMiddleware("delete_rolebinding")(rbac.RoleBindingsHandler(clientset)))
	api.GET("/rolebinding/details", middleware.RBACMiddleware("view_rolebinding_details")(rbac.RoleBindingDetailsHandler(clientset)))

	// Cluster role routes
	api.GET("/clusterroles", middleware.RBACMiddleware("list_clusterroles")(rbac.ClusterRolesHandler(clientset)))
	api.POST("/clusterroles", middleware.RBACMiddleware("create_clusterrole")(rbac.ClusterRolesHandler(clientset)))
	api.PUT("/clusterroles", middleware.RBACMiddleware("update_clusterrole")(rbac.ClusterRolesHandler(clientset)))
	api.DELETE("/clusterroles", middleware.RBACMiddleware("delete_clusterrole")(rbac.ClusterRolesHandler(clientset)))
	api.GET("/clusterroles/details", middleware.RBACMiddleware("view_clusterrole_details")(rbac.ClusterRoleDetailsHandler(clientset)))

	// Cluster role bindings routes
	api.GET("/clusterrolebindings", middleware.RBACMiddleware("list_clusterrolebindings")(rbac.ClusterRoleBindingsHandler(clientset)))
	api.POST("/clusterrolebindings", middleware.RBACMiddleware("create_clusterrolebinding")(rbac.ClusterRoleBindingsHandler(clientset)))
	api.PUT("/clusterrolebindings", middleware.RBACMiddleware("update_clusterrolebinding")(rbac.ClusterRoleBindingsHandler(clientset)))
	api.DELETE("/clusterrolebindings", middleware.RBACMiddleware("delete_clusterrolebinding")(rbac.ClusterRoleBindingsHandler(clientset)))
	api.GET("/clusterrolebinding/details", middleware.RBACMiddleware("view_clusterrolebinding_details")(rbac.ClusterRoleBindingDetailsHandler(clientset)))

	// Resource routes
	api.GET("/resources", middleware.RBACMiddleware("list_resources")(rbac.APIResourcesHandler(clientset)))

	// Account routes
	api.GET("/serviceaccounts", middleware.RBACMiddleware("list_serviceaccounts")(rbac.ServiceAccountsHandler(clientset)))
	api.POST("/serviceaccounts", middleware.RBACMiddleware("create_serviceaccount")(rbac.ServiceAccountsHandler(clientset)))
	api.DELETE("/serviceaccounts", middleware.RBACMiddleware("delete_serviceaccount")(rbac.ServiceAccountsHandler(clientset)))
	api.GET("/serviceaccount-details", middleware.RBACMiddleware("view_serviceaccount_details")(rbac.ServiceAccountDetailsHandler(clientset)))

	// User routes
	api.GET("/users", middleware.RBACMiddleware("list_users")(rbac.UsersHandler(clientset)))
	api.GET("/user-details", middleware.RBACMiddleware("view_user_details")(rbac.UserDetailsHandler(clientset)))

	// Group routes
	api.GET("/groups", middleware.RBACMiddleware("list_groups")(rbac.GroupsHandler(clientset)))
	api.GET("/group-details", middleware.RBACMiddleware("view_group_details")(rbac.GroupDetailsHandler(clientset)))

	// User roles route
	api.GET("/user-roles", middleware.RBACMiddleware("view_user_roles")(rbac.UserRolesHandler(clientset)))

	// Audit logs route
	api.GET("/audit-logs", middleware.RBACMiddleware("view_audit_logs")(handlers.GetAuditLogsHandler))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Root URL handler
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Welcome to the RBAC Manager"})
	})

	// Simulation route
	api.POST("/simulate", middleware.RBACMiddleware("simulate")(rbac.SimulateHandler(clientset)))
}