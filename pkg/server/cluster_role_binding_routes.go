package server

import (
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

func registerClusterRoleBindingRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	api.GET("/clusterrolebindings", middleware.AuthAndRBACMiddleware("list_clusterrolebindings", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.POST("/clusterrolebindings", middleware.AuthAndRBACMiddleware("create_clusterrolebinding", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.PUT("/clusterrolebindings", middleware.AuthAndRBACMiddleware("update_clusterrolebinding", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.DELETE("/clusterrolebindings", middleware.AuthAndRBACMiddleware("delete_clusterrolebinding", config.IsDevMode)(rbac.ClusterRoleBindingsHandler(clientset)))
	api.GET("/clusterrolebinding/details", middleware.AuthAndRBACMiddleware("view_clusterrolebinding_details", config.IsDevMode)(rbac.ClusterRoleBindingDetailsHandler(clientset)))
}
