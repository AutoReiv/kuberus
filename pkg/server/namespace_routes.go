package server

import (
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

func registerNamespaceRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	api.GET("/namespaces", middleware.AuthAndRBACMiddleware("list_namespaces", config.IsDevMode)(rbac.NamespacesHandler(clientset)))
	api.POST("/namespaces", middleware.AuthAndRBACMiddleware("create_namespace", config.IsDevMode)(rbac.NamespacesHandler(clientset)))
	api.DELETE("/namespaces", middleware.AuthAndRBACMiddleware("delete_namespace", config.IsDevMode)(rbac.NamespacesHandler(clientset)))
}
