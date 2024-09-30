package server

import (
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

func registerRoleBindingRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	api.GET("/rolebindings", middleware.AuthAndRBACMiddleware("list_rolebindings", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.POST("/rolebindings", middleware.AuthAndRBACMiddleware("create_rolebinding", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.PUT("/rolebindings", middleware.AuthAndRBACMiddleware("update_rolebinding", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.DELETE("/rolebindings", middleware.AuthAndRBACMiddleware("delete_rolebinding", config.IsDevMode)(rbac.RoleBindingsHandler(clientset)))
	api.GET("/rolebinding/details", middleware.AuthAndRBACMiddleware("view_rolebinding_details", config.IsDevMode)(rbac.RoleBindingDetailsHandler(clientset)))
}
