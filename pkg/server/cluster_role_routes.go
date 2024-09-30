package server

import (
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

func registerClusterRoleRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	api.GET("/clusterroles", middleware.AuthAndRBACMiddleware("list_clusterroles", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.POST("/clusterroles", middleware.AuthAndRBACMiddleware("create_clusterrole", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.PUT("/clusterroles", middleware.AuthAndRBACMiddleware("update_clusterrole", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.DELETE("/clusterroles", middleware.AuthAndRBACMiddleware("delete_clusterrole", config.IsDevMode)(rbac.ClusterRolesHandler(clientset)))
	api.GET("/clusterroles/details", middleware.AuthAndRBACMiddleware("view_clusterrole_details", config.IsDevMode)(rbac.ClusterRoleDetailsHandler(clientset)))
}
