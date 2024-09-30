package server

import (
	"rbac/pkg/handlers/rbac"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

func registerServiceAccountRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	api.GET("/serviceaccounts", middleware.AuthAndRBACMiddleware("list_serviceaccounts", config.IsDevMode)(rbac.ServiceAccountsHandler(clientset)))
	api.POST("/serviceaccounts", middleware.AuthAndRBACMiddleware("create_serviceaccount", config.IsDevMode)(rbac.ServiceAccountsHandler(clientset)))
	api.DELETE("/serviceaccounts", middleware.AuthAndRBACMiddleware("delete_serviceaccount", config.IsDevMode)(rbac.ServiceAccountsHandler(clientset)))
	api.GET("/serviceaccount-details", middleware.AuthAndRBACMiddleware("view_serviceaccount_details", config.IsDevMode)(rbac.ServiceAccountDetailsHandler(clientset)))
}
