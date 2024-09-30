package server

import (
	"rbac/pkg/handlers"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

func registerAdminRoutes(e *echo.Echo, clientset *kubernetes.Clientset, config *Config) {
	e.POST("/admin/create", handlers.CreateAdminHandler)
	e.POST("/admin/oidc/config", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(handlers.SetOIDCConfigHandler))
	uploadCertsHandler := handlers.NewUploadCertsHandler(clientset)
	e.POST("/admin/upload-certs", middleware.AuthAndRBACMiddleware("", config.IsDevMode)(uploadCertsHandler.ServeHTTP))
}