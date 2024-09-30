package server

import (
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

func registerRoleRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	api.GET("/roles", middleware.AuthAndRBACMiddleware("list_roles", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.POST("/roles", middleware.AuthAndRBACMiddleware("create_role", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.PUT("/roles", middleware.AuthAndRBACMiddleware("update_role", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.DELETE("/roles", middleware.AuthAndRBACMiddleware("delete_role", config.IsDevMode)(rbac.RolesHandler(clientset)))
	api.GET("/roles/details", middleware.AuthAndRBACMiddleware("view_role_details", config.IsDevMode)(rbac.RoleDetailsHandler(clientset)))
}
